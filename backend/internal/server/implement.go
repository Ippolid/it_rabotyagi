package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
	"it_rabotyagi/api/openapi"
	"it_rabotyagi/internal/business/models"
	"it_rabotyagi/internal/business/services"
	"it_rabotyagi/internal/data/repositories"
)

// PermissionChecker defines the interface to verify if a user may edit mutable resources.
type PermissionChecker interface {
	IsEditor(ctx context.Context, userID int) (bool, error)
}

// ServerImplementation реализует интерфейс openapi.ServerInterface
type ServerImplementation struct {
	authService    *services.AuthService
	repo           *repositories.UserRepository
	sessionRepo    *repositories.SessionRepository
	questionRepo   *repositories.QuestionRepository
	courseRepo     *repositories.CourseRepository
	mentorRepo     *repositories.MentorRepository
	permissionRepo PermissionChecker
	statisticsRepo *repositories.StatisticsRepository
	userService    *services.UserService
}

func NewServerImplementation(authService *services.AuthService, repo *repositories.UserRepository, sessionRepo *repositories.SessionRepository, questionRepo *repositories.QuestionRepository, courseRepo *repositories.CourseRepository, mentorRepo *repositories.MentorRepository, permissionRepo PermissionChecker, statisticsRepo *repositories.StatisticsRepository, userService *services.UserService) *ServerImplementation {
	return &ServerImplementation{
		authService:    authService,
		repo:           repo,
		sessionRepo:    sessionRepo,
		questionRepo:   questionRepo,
		courseRepo:     courseRepo,
		mentorRepo:     mentorRepo,
		permissionRepo: permissionRepo,
		statisticsRepo: statisticsRepo,
		userService:    userService,
	}
}

// HashPassword хеширует пароль используя bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword проверяет соответствие пароля хешу
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *ServerImplementation) ensureEditor(ctx echo.Context) bool {
	userID, ok := GetUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
		return false
	}
	role, ok := GetRole(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
		return false
	}
	if role == "admin" {
		return true
	}
	allowed, err := s.permissionRepo.IsEditor(ctx.Request().Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to check permissions",
			Code:    strPtr("PERMISSION_CHECK_ERROR"),
		})
		return false
	}
	if !allowed {
		ctx.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Message: "Insufficient permissions",
			Code:    strPtr("FORBIDDEN"),
		})
		return false
	}
	return true
}

// RegisterUser регистрирует нового пользователя
// (POST /auth/register)
func (s *ServerImplementation) RegisterUser(ctx echo.Context) error {
	var req openapi.AuthRegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Проверяем существование пользователя по email и username
	existingUserID, err := s.repo.CheckUser(ctx.Request().Context(), string(req.Email), req.Nickname)
	if err != nil {
		// Реальная ошибка БД
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to check user",
			Code:    strPtr("USER_CHECK_ERROR"),
		})
	}

	if existingUserID != nil {
		// Пользователь уже существует
		return ctx.JSON(http.StatusConflict, openapi.ErrorResponse{
			Message: "User with this email or username already exists",
			Code:    strPtr("USER_ALREADY_EXISTS"),
		})
	}

	// Хешируем пароль
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to hash password",
			Code:    strPtr("PASSWORD_HASH_ERROR"),
		})
	}

	userID, err := s.repo.CreateUser(ctx.Request().Context(), string(req.Email), req.Nickname, hashedPassword)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to create user",
			Code:    strPtr("USER_CREATION_ERROR"),
		})
	}

	accessToken, refreshToken, expiresIn, err := s.authService.GenerateTokens(userID, req.Nickname, "user")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to generate tokens",
			Code:    strPtr("TOKEN_GENERATION_ERROR"),
		})
	}

	refreshExpiration := s.authService.GetRefreshTokenExpiration()
	err = s.sessionRepo.CreateSession(ctx.Request().Context(), userID, refreshToken,
		time.Now().Add(refreshExpiration))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to create session",
			Code:    strPtr("SESSION_CREATION_ERROR"),
		})
	}

	tokens := openapi.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	return ctx.JSON(http.StatusCreated, tokens)
}

// LoginUser выполняет вход пользователя
// (POST /auth/login)
func (s *ServerImplementation) LoginUser(ctx echo.Context) error {
	var req openapi.AuthLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Получаем пользователя из БД по email
	user, err := s.repo.GetUserByEmail(ctx.Request().Context(), string(req.Email))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Invalid email or password",
			Code:    strPtr("INVALID_CREDENTIALS"),
		})
	}

	// Проверяем пароль
	if err := VerifyPassword(user.Password, req.Password); err != nil {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Invalid email or password",
			Code:    strPtr("INVALID_CREDENTIALS"),
		})
	}

	// Генерируем токены
	accessToken, refreshToken, expiresIn, err := s.authService.GenerateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to generate tokens",
			Code:    strPtr("TOKEN_GENERATION_ERROR"),
		})
	}

	// Сохраняем сессию в БД
	refreshExpiration := s.authService.GetRefreshTokenExpiration()
	err = s.sessionRepo.CreateSession(ctx.Request().Context(), user.ID, refreshToken,
		time.Now().Add(refreshExpiration))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to create session",
			Code:    strPtr("SESSION_CREATION_ERROR"),
		})
	}

	tokens := openapi.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	return ctx.JSON(http.StatusOK, tokens)
}

// RefreshTokens обновляет access токен
// (POST /auth/refresh)
func (s *ServerImplementation) RefreshTokens(ctx echo.Context) error {
	var req openapi.AuthRefreshRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Проверяем наличие сессии в БД
	session, err := s.sessionRepo.GetSessionByToken(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Invalid or revoked refresh token",
			Code:    strPtr("INVALID_REFRESH_TOKEN"),
		})
	}

	// Проверяем, не истекла ли сессия
	if session.ExpiresAt.Before(time.Now()) {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Refresh token expired",
			Code:    strPtr("REFRESH_TOKEN_EXPIRED"),
		})
	}

	// Валидируем токен через AuthService
	accessToken, expiresIn, err := s.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Invalid refresh token",
			Code:    strPtr("INVALID_REFRESH_TOKEN"),
		})
	}

	tokens := openapi.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		ExpiresIn:    expiresIn,
	}

	return ctx.JSON(http.StatusOK, tokens)
}

// GetCurrentUser получает профиль текущего пользователя
// (GET /users/me)
func (s *ServerImplementation) GetCurrentUser(ctx echo.Context) error {
	// Извлекаем данные пользователя из контекста (добавлены middleware)
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "User not authenticated",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	role, _ := GetRole(ctx)

	// Получить полную информацию о пользователе из БД по userID
	user, err := s.repo.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to get user info",
			Code:    strPtr("USER_FETCH_ERROR"),
		})
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	fullName := user.Username
	if user.Name != nil {
		fullName = *user.Name
	}

	profile := openapi.UserProfile{
		Id:          user.ID,
		Email:       openapi_types.Email(email),
		FullName:    fullName,
		Role:        role,
		AvatarUrl:   user.AvatarURL,
		Description: user.Description,
	}

	return ctx.JSON(http.StatusOK, profile)
}

// ListMentors получает список менторов
// (GET /mentors)
func (s *ServerImplementation) ListMentors(ctx echo.Context, params openapi.ListMentorsParams) error {
	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	mentors, total, err := s.mentorRepo.ListMentors(ctx.Request().Context(), params.Specialization, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch mentors",
			Code:    strPtr("MENTORS_FETCH_ERROR"),
		})
	}

	items := make([]openapi.MentorCard, 0, len(mentors))
	for _, m := range mentors {
		items = append(items, openapi.MentorCard{
			Id:                m.ID,
			FullName:          m.Specialization,
			Title:             m.Specialization,
			Skills:            m.Tags,
			YearsOfExperience: m.ExperienceYears,
		})
	}

	totalCopy := total
	return ctx.JSON(http.StatusOK, openapi.MentorList{Items: items, Total: &totalCopy})
}

// GetMentorById возвращает профиль ментора
// (GET /mentors/{id})
func (s *ServerImplementation) GetMentorById(ctx echo.Context, id int) error {
	mentor, err := s.mentorRepo.GetMentor(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Mentor not found",
			Code:    strPtr("MENTOR_NOT_FOUND"),
		})
	}

	resp := openapi.MentorProfile{
		Id:              mentor.ID,
		UserId:          mentor.UserID,
		Specialization:  mentor.Specialization,
		Grade:           mentor.Grade,
		ExperienceYears: mentor.ExperienceYears,
		Description:     mentor.Description,
		Pricelist:       nil,
		Contacts:        nil,
	}
	if len(mentor.Tags) > 0 {
		resp.Tags = &mentor.Tags
	}
	if mentor.Contacts != nil {
		resp.Contacts = &mentor.Contacts
	}
	if mentor.Pricelist != nil {
		resp.Pricelist = &mentor.Pricelist
	}

	return ctx.JSON(http.StatusOK, resp)
}

// Logout отзывает текущую сессию пользователя
// (POST /auth/logout)
func (s *ServerImplementation) Logout(ctx echo.Context) error {
	var req openapi.AuthRefreshRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Отзываем сессию
	err := s.sessionRepo.RevokeSession(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to revoke session",
			Code:    strPtr("SESSION_REVOKE_ERROR"),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

// ListQuestions получает список всех вопросов
// (GET /questions)
func (s *ServerImplementation) ListQuestions(ctx echo.Context, params openapi.ListQuestionsParams) error {
	// Определяем лимит и оффсет
	limit := 20 // по умолчанию
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	// Формируем фильтры
	filters := repositories.QuestionFilters{
		Search:     params.Search,
		Technology: params.Technology,
		Company:    params.Company,
	}

	// Преобразуем difficulty из enum в string
	if params.Difficulty != nil {
		difficultyStr := string(*params.Difficulty)
		filters.Difficulty = &difficultyStr
	}

	// Получаем вопросы из БД
	questions, total, err := s.questionRepo.GetAllQuestions(ctx.Request().Context(), filters, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch questions",
			Code:    strPtr("QUESTIONS_FETCH_ERROR"),
		})
	}

	// Получаем список решенных вопросов для авторизованного пользователя
	var solvedMap map[int]bool
	if userID, ok := GetUserID(ctx); ok {
		solvedMap, err = s.statisticsRepo.GetSolvedQuestionIDs(ctx.Request().Context(), userID)
		if err != nil {
			// Не критичная ошибка, просто логируем и продолжаем без маркировки
			solvedMap = make(map[int]bool)
		}
	} else {
		solvedMap = make(map[int]bool)
	}

	// Преобразуем в формат OpenAPI
	items := make([]openapi.QuestionListItem, 0, len(questions))
	for _, q := range questions {
		companyTags := make([]string, 0)
		if q.CompanyTags != nil {
			companyTags = q.CompanyTags
		}

		solved := solvedMap[q.ID]

		items = append(items, openapi.QuestionListItem{
			Id:          q.ID,
			Title:       q.Title,
			Technology:  q.Technology,
			Difficulty:  openapi.QuestionListItemDifficulty(q.Difficulty),
			CompanyTags: &companyTags,
			Solved:      &solved,
		})
	}

	questionList := openapi.QuestionList{
		Items: items,
		Total: &total,
	}

	return ctx.JSON(http.StatusOK, questionList)
}

// GetQuestionById получает полную информацию о вопросе по ID
// (GET /questions/{id})
func (s *ServerImplementation) GetQuestionById(ctx echo.Context, id int) error {
	question, err := s.questionRepo.GetQuestionByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Question not found",
			Code:    strPtr("QUESTION_NOT_FOUND"),
		})
	}

	resp := openapi.QuestionDetail{
		Id:         question.ID,
		Title:      question.Title,
		Content:    question.Content,
		Difficulty: openapi.QuestionDetailDifficulty(question.Difficulty),
		Technology: question.Technology,
	}

	// Опциональные поля
	if len(question.Options) > 0 {
		resp.Options = &question.Options
	}
	if question.CorrectAnswer != "" {
		resp.CorrectAnswer = &question.CorrectAnswer
	}
	if question.Explanation != nil {
		resp.Explanation = question.Explanation
	}
	if len(question.Tags) > 0 {
		resp.Tags = &question.Tags
	}

	// Проверяем, является ли вопрос открытым
	isOpenEnded := len(question.Options) == 0 || question.CorrectAnswer == ""
	resp.IsOpenEnded = &isOpenEnded

	return ctx.JSON(http.StatusOK, resp)
}

// CreateQuestion создает новый вопрос
// (POST /questions)
func (s *ServerImplementation) CreateQuestion(ctx echo.Context) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	var req openapi.QuestionCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	tags := []string{}
	if req.Tags != nil {
		tags = *req.Tags
	}

	options := []string{}
	if req.Options != nil {
		options = *req.Options
	}

	correctAnswer := ""
	if req.CorrectAnswer != nil {
		correctAnswer = *req.CorrectAnswer
	}

	id, err := s.questionRepo.CreateQuestion(ctx.Request().Context(), req.Title, req.Content, string(req.Difficulty), req.Technology, tags, options, correctAnswer, req.Explanation)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to create question",
			Code:    strPtr("QUESTION_CREATE_ERROR"),
		})
	}

	return s.GetQuestionById(ctx, id)
}

// UpdateQuestion обновляет существующий вопрос
// (PATCH /questions/{id})
func (s *ServerImplementation) UpdateQuestion(ctx echo.Context, id int) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	var req openapi.QuestionUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	var difficulty *string
	if req.Difficulty != nil {
		val := string(*req.Difficulty)
		difficulty = &val
	}

	if err := s.questionRepo.UpdateQuestion(ctx.Request().Context(), id, req.Title, req.Content, difficulty, req.Technology, req.Tags, req.Options, req.CorrectAnswer, req.Explanation); err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Message: "Question not found",
				Code:    strPtr("QUESTION_NOT_FOUND"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to update question",
			Code:    strPtr("QUESTION_UPDATE_ERROR"),
		})
	}

	return s.GetQuestionById(ctx, id)
}

// DeleteQuestion удаляет вопрос
// (DELETE /questions/{id})
func (s *ServerImplementation) DeleteQuestion(ctx echo.Context, id int) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	if err := s.questionRepo.DeleteQuestion(ctx.Request().Context(), id); err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Message: "Question not found",
				Code:    strPtr("QUESTION_NOT_FOUND"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to delete question",
			Code:    strPtr("QUESTION_DELETE_ERROR"),
		})
	}
	return ctx.NoContent(http.StatusNoContent)
}

// SubmitQuestionAnswer принимает ответ пользователя на вопрос
// (POST /questions/{id}/submit)
func (s *ServerImplementation) SubmitQuestionAnswer(ctx echo.Context, id int) error {
	// Получаем userID из контекста
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "User not authenticated",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	// Парсим request body
	var req openapi.QuestionSubmitRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Получаем вопрос из БД
	question, err := s.questionRepo.GetQuestionByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Question not found",
			Code:    strPtr("QUESTION_NOT_FOUND"),
		})
	}

	// Проверяем правильность ответа
	isCorrect := req.UserAnswer == question.CorrectAnswer

	// Проверяем, был ли вопрос уже решен
	_, alreadySolved, err := s.statisticsRepo.CheckIfQuestionSolved(ctx.Request().Context(), userID, id)
	if err != nil {
		// Логируем ошибку, но не останавливаем выполнение
		fmt.Printf("Error checking if question solved: %v\n", err)
	}

	// Сохраняем попытку (предполагаем, что пользователь потратил 30 секунд)
	timeSpent := 30 * time.Second
	if err := s.statisticsRepo.RecordQuestionAttempt(ctx.Request().Context(), userID, id, isCorrect, timeSpent); err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to record answer",
			Code:    strPtr("RECORD_ANSWER_ERROR"),
		})
	}

	// Формируем ответ
	resp := openapi.QuestionSubmitResponse{
		IsCorrect:     isCorrect,
		CorrectAnswer: question.CorrectAnswer,
		Explanation:   question.Explanation,
		AlreadySolved: &alreadySolved,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// ListCourses получает список курсов
// (GET /courses)
func (s *ServerImplementation) ListCourses(ctx echo.Context, params openapi.ListCoursesParams) error {
	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	courses, total, err := s.courseRepo.ListCourses(ctx.Request().Context(), limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch courses",
			Code:    strPtr("COURSES_FETCH_ERROR"),
		})
	}

	items := make([]openapi.CourseListItem, 0, len(courses))
	for _, c := range courses {
		items = append(items, openapi.CourseListItem{
			Id:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			IsPublished: c.IsPublished,
		})
	}
	totalCopy := total
	return ctx.JSON(http.StatusOK, openapi.CourseList{Items: items, Total: &totalCopy})
}

// CreateCourse создает новый курс
// (POST /courses)
func (s *ServerImplementation) CreateCourse(ctx echo.Context) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	var req openapi.CourseCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	isPublished := false
	if req.IsPublished != nil {
		isPublished = *req.IsPublished
	}

	id, err := s.courseRepo.CreateCourse(ctx.Request().Context(), req.Title, req.Description, isPublished)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to create course",
			Code:    strPtr("COURSE_CREATE_ERROR"),
		})
	}

	return s.GetCourseById(ctx, id)
}

// GetCourseById возвращает курс
// (GET /courses/{id})
func (s *ServerImplementation) GetCourseById(ctx echo.Context, id int) error {
	course, err := s.courseRepo.GetCourse(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Course not found",
			Code:    strPtr("COURSE_NOT_FOUND"),
		})
	}

	modules, err := s.courseRepo.ListModules(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch modules",
			Code:    strPtr("MODULES_FETCH_ERROR"),
		})
	}

	moduleItems := make([]openapi.ModuleDetail, 0, len(modules))
	for _, m := range modules {
		moduleItems = append(moduleItems, openapi.ModuleDetail{
			Id:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Content:     m.Content,
			Order:       m.Order,
		})
	}

	resp := openapi.CourseDetail{
		Id:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		IsPublished: course.IsPublished,
		Modules:     moduleItems,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// UpdateCourse обновляет курс
// (PATCH /courses/{id})
func (s *ServerImplementation) UpdateCourse(ctx echo.Context, id int) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	var req openapi.CourseUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	if err := s.courseRepo.UpdateCourse(ctx.Request().Context(), id, req.Title, req.Description, req.IsPublished); err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Message: "Course not found",
				Code:    strPtr("COURSE_NOT_FOUND"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to update course",
			Code:    strPtr("COURSE_UPDATE_ERROR"),
		})
	}

	return s.GetCourseById(ctx, id)
}

// DeleteCourse удаляет курс
// (DELETE /courses/{id})
func (s *ServerImplementation) DeleteCourse(ctx echo.Context, id int) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	if err := s.courseRepo.DeleteCourse(ctx.Request().Context(), id); err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Message: "Course not found",
				Code:    strPtr("COURSE_NOT_FOUND"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to delete course",
			Code:    strPtr("COURSE_DELETE_ERROR"),
		})
	}
	return ctx.NoContent(http.StatusNoContent)
}

// ListModules возвращает модули курса
// (GET /courses/{id}/modules)
func (s *ServerImplementation) ListModules(ctx echo.Context, id int) error {
	if _, err := s.courseRepo.GetCourse(ctx.Request().Context(), id); err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Course not found",
			Code:    strPtr("COURSE_NOT_FOUND"),
		})
	}

	modules, err := s.courseRepo.ListModules(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch modules",
			Code:    strPtr("MODULES_FETCH_ERROR"),
		})
	}

	items := make([]openapi.ModuleDetail, 0, len(modules))
	for _, m := range modules {
		items = append(items, openapi.ModuleDetail{
			Id:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Content:     m.Content,
			Order:       m.Order,
		})
	}
	return ctx.JSON(http.StatusOK, openapi.ModuleList{Items: items})
}

// CreateModule добавляет модуль
// (POST /courses/{id}/modules)
func (s *ServerImplementation) CreateModule(ctx echo.Context, id int) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	var req openapi.ModuleCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	if _, err := s.courseRepo.GetCourse(ctx.Request().Context(), id); err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Course not found",
			Code:    strPtr("COURSE_NOT_FOUND"),
		})
	}

	moduleID, err := s.courseRepo.CreateModule(ctx.Request().Context(), id, req.Title, req.Description, req.Order, req.Content)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to create module",
			Code:    strPtr("MODULE_CREATE_ERROR"),
		})
	}

	return s.GetModuleById(ctx, id, moduleID)
}

// GetModuleById возвращает модуль
// (GET /courses/{id}/modules/{moduleId})
func (s *ServerImplementation) GetModuleById(ctx echo.Context, id int, moduleId int) error {
	module, err := s.courseRepo.GetModule(ctx.Request().Context(), id, moduleId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Module not found",
			Code:    strPtr("MODULE_NOT_FOUND"),
		})
	}

	resp := openapi.ModuleDetail{
		Id:          module.ID,
		Title:       module.Title,
		Description: module.Description,
		Content:     module.Content,
		Order:       module.Order,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// UpdateModule обновляет модуль
// (PATCH /courses/{id}/modules/{moduleId})
func (s *ServerImplementation) UpdateModule(ctx echo.Context, id int, moduleId int) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	var req openapi.ModuleUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	if err := s.courseRepo.UpdateModule(ctx.Request().Context(), id, moduleId, req.Title, req.Description, req.Content, req.Order); err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Message: "Module not found",
				Code:    strPtr("MODULE_NOT_FOUND"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to update module",
			Code:    strPtr("MODULE_UPDATE_ERROR"),
		})
	}

	return s.GetModuleById(ctx, id, moduleId)
}

// DeleteModule удаляет модуль
// (DELETE /courses/{id}/modules/{moduleId})
func (s *ServerImplementation) DeleteModule(ctx echo.Context, id int, moduleId int) error {
	if !s.ensureEditor(ctx) {
		return nil
	}

	if err := s.courseRepo.DeleteModule(ctx.Request().Context(), id, moduleId); err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Message: "Module not found",
				Code:    strPtr("MODULE_NOT_FOUND"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to delete module",
			Code:    strPtr("MODULE_DELETE_ERROR"),
		})
	}
	return ctx.NoContent(http.StatusNoContent)
}

// EnrollCourse записывает пользователя на курс
// (POST /courses/{id}/enroll)
func (s *ServerImplementation) EnrollCourse(ctx echo.Context, id int) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	if _, err := s.courseRepo.GetCourse(ctx.Request().Context(), id); err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Course not found",
			Code:    strPtr("COURSE_NOT_FOUND"),
		})
	}

	modules, err := s.courseRepo.ListModules(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch modules",
			Code:    strPtr("MODULES_FETCH_ERROR"),
		})
	}

	startedAt, err := s.courseRepo.EnrollUser(ctx.Request().Context(), userID, id, len(modules))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to enroll",
			Code:    strPtr("ENROLL_ERROR"),
		})
	}

	return ctx.JSON(http.StatusOK, openapi.EnrollmentStatus{
		CourseId:   id,
		Enrolled:   true,
		EnrolledAt: &startedAt,
	})
}

// GetCourseProgress возвращает прогресс
// (GET /courses/{id}/progress)
func (s *ServerImplementation) GetCourseProgress(ctx echo.Context, id int) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	if _, err := s.courseRepo.GetCourse(ctx.Request().Context(), id); err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Course not found",
			Code:    strPtr("COURSE_NOT_FOUND"),
		})
	}

	total, completed, modules, err := s.courseRepo.GetCourseProgress(ctx.Request().Context(), userID, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch progress",
			Code:    strPtr("PROGRESS_FETCH_ERROR"),
		})
	}

	moduleItems := make([]openapi.ModuleProgressItem, 0, len(modules))
	for _, m := range modules {
		moduleItems = append(moduleItems, openapi.ModuleProgressItem{
			ModuleId:    m.Module.ID,
			Title:       m.Module.Title,
			Completed:   m.Completed,
			CompletedAt: m.CompletedAt,
		})
	}

	var percentage float32
	if total > 0 {
		percentage = float32(completed) * 100.0 / float32(total)
	}

	return ctx.JSON(http.StatusOK, openapi.CourseProgress{
		CourseId:         id,
		TotalModules:     total,
		CompletedModules: completed,
		Percentage:       &percentage,
		Modules:          &moduleItems,
	})
}

// CompleteModule отмечает модуль завершенным
// (POST /courses/{courseId}/modules/{moduleId}/complete)
func (s *ServerImplementation) CompleteModule(ctx echo.Context, courseId int, moduleId int) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	enrolled, err := s.courseRepo.IsUserEnrolled(ctx.Request().Context(), userID, courseId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to check enrollment",
			Code:    strPtr("ENROLL_CHECK_ERROR"),
		})
	}
	if !enrolled {
		return ctx.JSON(http.StatusForbidden, openapi.ErrorResponse{
			Message: "Enroll to course first",
			Code:    strPtr("ENROLL_REQUIRED"),
		})
	}

	completedAt, err := s.courseRepo.CompleteModule(ctx.Request().Context(), userID, courseId, moduleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
				Message: "Module not found",
				Code:    strPtr("MODULE_NOT_FOUND"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to complete module",
			Code:    strPtr("MODULE_COMPLETE_ERROR"),
		})
	}

	return ctx.JSON(http.StatusOK, openapi.ModuleCompletion{
		CourseId:    courseId,
		ModuleId:    moduleId,
		Completed:   true,
		CompletedAt: completedAt,
	})
}

// GetUserStatistics возвращает общую статистику пользователя
// (GET /users/me/statistics)
func (s *ServerImplementation) GetUserStatistics(ctx echo.Context) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	stats, err := s.statisticsRepo.GetOverallStatistics(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to get statistics",
			Code:    strPtr("STATISTICS_FETCH_ERROR"),
		})
	}

	response := openapi.UserStatistics{
		CoursesEnrolled:    stats.CoursesEnrolled,
		CoursesCompleted:   stats.CoursesCompleted,
		OverallProgressPct: stats.OverallProgressPct,
		QuestionsStatistics: openapi.QuestionStatisticsSummary{
			TotalSolved:       stats.QuestionStatistics.TotalSolved,
			CorrectAnswersPct: stats.QuestionStatistics.CorrectAnswersPct,
			TotalAttempts:     stats.QuestionStatistics.TotalAttempts,
			AvgTimeSpent:      stats.QuestionStatistics.AvgTimeSpent,
		},
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetUserCourseStatistics возвращает детальную статистику по курсам
// (GET /users/me/statistics/courses)
func (s *ServerImplementation) GetUserCourseStatistics(ctx echo.Context) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	stats, err := s.statisticsRepo.GetCourseStatistics(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to get course statistics",
			Code:    strPtr("STATISTICS_FETCH_ERROR"),
		})
	}

	items := make([]openapi.UserCourseStatisticsItem, 0, len(stats))
	for _, stat := range stats {
		items = append(items, openapi.UserCourseStatisticsItem{
			CourseId:         stat.CourseID,
			CourseTitle:      stat.CourseTitle,
			TotalModules:     stat.TotalModules,
			CompletedModules: stat.CompletedModules,
			ProgressPct:      stat.ProgressPct,
			StartedAt:        stat.StartedAt,
			TimeSpent:        &stat.TimeSpent,
		})
	}

	return ctx.JSON(http.StatusOK, openapi.UserCourseStatisticsList{
		Items: items,
	})
}

// GetUserQuestionStatistics возвращает статистику по вопросам
// (GET /users/me/statistics/questions)
func (s *ServerImplementation) GetUserQuestionStatistics(ctx echo.Context, params openapi.GetUserQuestionStatisticsParams) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	overallStats, err := s.statisticsRepo.GetQuestionStatisticsOverall(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to get question statistics",
			Code:    strPtr("STATISTICS_FETCH_ERROR"),
		})
	}

	byCourseStats, err := s.statisticsRepo.GetQuestionStatisticsByCourse(ctx.Request().Context(), userID, params.CourseId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to get question statistics by course",
			Code:    strPtr("STATISTICS_FETCH_ERROR"),
		})
	}

	byCourseItems := make([]openapi.QuestionStatisticsByCourse, 0, len(byCourseStats))
	for _, stat := range byCourseStats {
		byCourseItems = append(byCourseItems, openapi.QuestionStatisticsByCourse{
			CourseId:    stat.CourseID,
			CourseTitle: stat.CourseTitle,
			Statistics: openapi.QuestionStatisticsSummary{
				TotalSolved:       stat.Statistics.TotalSolved,
				CorrectAnswersPct: stat.Statistics.CorrectAnswersPct,
				TotalAttempts:     stat.Statistics.TotalAttempts,
				AvgTimeSpent:      stat.Statistics.AvgTimeSpent,
			},
		})
	}

	response := openapi.UserQuestionStatistics{
		Overall: openapi.QuestionStatisticsSummary{
			TotalSolved:       overallStats.TotalSolved,
			CorrectAnswersPct: overallStats.CorrectAnswersPct,
			TotalAttempts:     overallStats.TotalAttempts,
			AvgTimeSpent:      overallStats.AvgTimeSpent,
		},
		ByCourse: byCourseItems,
	}

	return ctx.JSON(http.StatusOK, response)
}

// UpdateUserProfile обновляет профиль пользователя
// (PATCH /users/me/profile)
func (s *ServerImplementation) UpdateUserProfile(ctx echo.Context) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	var req openapi.UserProfileUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Конвертируем openapi типы в business модели
	var emailStr *string
	if req.Email != nil {
		email := string(*req.Email)
		emailStr = &email
	}

	update := &models.UserProfileUpdate{
		Username:    req.Username,
		Name:        req.Name,
		Email:       emailStr,
		Description: req.Description,
	}

	// Валидация через UserService
	if err := s.userService.ValidateProfileUpdate(update); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: err.Error(),
			Code:    strPtr("VALIDATION_ERROR"),
		})
	}

	// Проверка уникальности username
	usernameChanged := false
	if req.Username != nil {
		exists, err := s.repo.CheckUsernameExists(ctx.Request().Context(), *req.Username, userID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
				Message: "Failed to check username",
				Code:    strPtr("USERNAME_CHECK_ERROR"),
			})
		}
		if exists {
			return ctx.JSON(http.StatusConflict, openapi.ErrorResponse{
				Message: "Username already exists",
				Code:    strPtr("USERNAME_CONFLICT"),
			})
		}
		usernameChanged = true
	}

	// Проверка уникальности email
	if emailStr != nil {
		exists, err := s.repo.CheckEmailExists(ctx.Request().Context(), *emailStr, userID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
				Message: "Failed to check email",
				Code:    strPtr("EMAIL_CHECK_ERROR"),
			})
		}
		if exists {
			return ctx.JSON(http.StatusConflict, openapi.ErrorResponse{
				Message: "Email already exists",
				Code:    strPtr("EMAIL_CONFLICT"),
			})
		}
	}

	// Обновление профиля
	if err := s.repo.UpdateUserProfile(ctx.Request().Context(), userID, req.Username, req.Name, emailStr, req.Description); err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to update profile",
			Code:    strPtr("PROFILE_UPDATE_ERROR"),
		})
	}

	// Получаем обновленный профиль
	user, err := s.repo.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to get updated profile",
			Code:    strPtr("PROFILE_FETCH_ERROR"),
		})
	}

	// Безопасная обработка nullable полей
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	fullName := user.Username // Fallback to username
	if user.Name != nil {
		fullName = *user.Name
	}

	profile := openapi.UserProfile{
		Id:       user.ID,
		Email:    openapi_types.Email(email),
		FullName: fullName,
		Role:     user.Role,
	}

	response := openapi.UserProfileUpdateResponse{
		Profile: profile,
	}

	// Если username изменен, генерируем новые токены
	if usernameChanged && req.Username != nil {
		// Отзываем все существующие сессии пользователя
		if err := s.sessionRepo.RevokeAllUserSessions(ctx.Request().Context(), userID); err != nil {
			return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
				Message: "Failed to revoke old sessions",
				Code:    strPtr("SESSION_REVOKE_ERROR"),
			})
		}

		accessToken, refreshToken, expiresIn, err := s.authService.GenerateTokens(userID, *req.Username, user.Role)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
				Message: "Failed to generate tokens",
				Code:    strPtr("TOKEN_GENERATION_ERROR"),
			})
		}

		// Создаем новую сессию
		refreshExpiration := s.authService.GetRefreshTokenExpiration()
		err = s.sessionRepo.CreateSession(ctx.Request().Context(), userID, refreshToken, time.Now().Add(refreshExpiration))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
				Message: "Failed to create session",
				Code:    strPtr("SESSION_CREATION_ERROR"),
			})
		}

		tokens := &openapi.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expiresIn,
		}
		response.Tokens = tokens
	}

	return ctx.JSON(http.StatusOK, response)
}

// UpdateUserAvatar обновляет аватар пользователя
// (PATCH /users/me/avatar)
func (s *ServerImplementation) UpdateUserAvatar(ctx echo.Context) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	var req openapi.UserAvatarUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Валидация URL
	if err := s.userService.ValidateAvatarURL(req.AvatarUrl); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: err.Error(),
			Code:    strPtr("VALIDATION_ERROR"),
		})
	}

	// Обновление аватара
	if err := s.repo.UpdateUserAvatar(ctx.Request().Context(), userID, req.AvatarUrl); err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to update avatar",
			Code:    strPtr("AVATAR_UPDATE_ERROR"),
		})
	}

	return ctx.JSON(http.StatusOK, openapi.UserAvatarResponse{
		AvatarUrl: req.AvatarUrl,
	})
}

// ChangeUserPassword изменяет пароль пользователя
// (POST /users/me/password)
func (s *ServerImplementation) ChangeUserPassword(ctx echo.Context) error {
	userID, ok := GetUserID(ctx)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Unauthorized",
			Code:    strPtr("UNAUTHORIZED"),
		})
	}

	var req openapi.ChangePasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: "Invalid request body",
			Code:    strPtr("INVALID_REQUEST"),
		})
	}

	// Валидация нового пароля
	if err := s.userService.ValidatePassword(req.NewPassword); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ErrorResponse{
			Message: err.Error(),
			Code:    strPtr("VALIDATION_ERROR"),
		})
	}

	// Получаем текущего пользователя для проверки старого пароля
	user, err := s.repo.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to get user",
			Code:    strPtr("USER_FETCH_ERROR"),
		})
	}

	// Проверяем старый пароль
	if err := VerifyPassword(user.Password, req.OldPassword); err != nil {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Invalid old password",
			Code:    strPtr("INVALID_OLD_PASSWORD"),
		})
	}

	// Хешируем новый пароль
	newPasswordHash, err := HashPassword(req.NewPassword)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to hash new password",
			Code:    strPtr("PASSWORD_HASH_ERROR"),
		})
	}

	// Обновляем пароль
	if err := s.repo.UpdateUserPassword(ctx.Request().Context(), userID, newPasswordHash); err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to update password",
			Code:    strPtr("PASSWORD_UPDATE_ERROR"),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password successfully changed",
	})
}

// strPtr возвращает указатель на строку
func strPtr(s string) *string {
	return &s
}
