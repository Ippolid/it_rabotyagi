package server

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"it_rabotyagi/api/openapi"
	"it_rabotyagi/internal/business/services"
	"it_rabotyagi/internal/data/repositories"
	"net/http"
)

// ServerImplementation реализует интерфейс openapi.ServerInterface
type ServerImplementation struct {
	authService  *services.AuthService
	repo         *repositories.UserRepository
	sessionRepo  *repositories.SessionRepository
	questionRepo *repositories.QuestionRepository
}

func NewServerImplementation(authService *services.AuthService, repo *repositories.UserRepository, sessionRepo *repositories.SessionRepository, questionRepo *repositories.QuestionRepository) *ServerImplementation {
	return &ServerImplementation{
		authService:  authService,
		repo:         repo,
		sessionRepo:  sessionRepo,
		questionRepo: questionRepo,
	}
}

func Hash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
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

	userID, err := s.repo.CreateUser(ctx.Request().Context(), string(req.Email), req.Nickname, Hash(req.Password))
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
	err = s.sessionRepo.CreateSession(ctx.Request().Context(), int64(userID), refreshToken,
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
	if Hash(req.Password) != user.Password {
		return ctx.JSON(http.StatusUnauthorized, openapi.ErrorResponse{
			Message: "Invalid email or password",
			Code:    strPtr("INVALID_CREDENTIALS"),
		})
	}

	// Генерируем токены
	accessToken, refreshToken, expiresIn, err := s.authService.GenerateTokens(int(user.ID), user.Username, user.Role)
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
		Id:       user.ID,
		Email:    openapi_types.Email(email),
		FullName: fullName,
		Role:     role,
	}

	return ctx.JSON(http.StatusOK, profile)
}

// ListMentors получает список менторов
// (GET /mentors)
func (s *ServerImplementation) ListMentors(ctx echo.Context, params openapi.ListMentorsParams) error {
	// TODO: Реализовать логику получения списка менторов
	mentorList := openapi.MentorList{
		Items: []openapi.MentorCard{
			{
				Id:       1,
				FullName: "Alice Smith",
				Title:    "Senior Go Developer",
				Skills:   []string{"Go", "Docker", "Kubernetes"},
				YearsOfExperience: func() *int {
					i := 5
					return &i
				}(),
			},
			{
				Id:       2,
				FullName: "Bob Johnson",
				Title:    "Backend Engineer",
				Skills:   []string{"Go", "PostgreSQL", "gRPC"},
				YearsOfExperience: func() *int {
					i := 7
					return &i
				}(),
			},
		},
		Total: func() *int {
			i := 2
			return &i
		}(),
	}

	return ctx.JSON(http.StatusOK, mentorList)
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
	limit := 30 // по умолчанию
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	// Получаем вопросы из БД
	questions, total, err := s.questionRepo.GetAllQuestions(ctx.Request().Context(), params.Technology, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, openapi.ErrorResponse{
			Message: "Failed to fetch questions",
			Code:    strPtr("QUESTIONS_FETCH_ERROR"),
		})
	}

	// Преобразуем в формат OpenAPI
	items := make([]openapi.QuestionListItem, 0, len(questions))
	for _, q := range questions {
		items = append(items, openapi.QuestionListItem{
			Id:         q.ID,
			Title:      q.Title,
			Technology: q.Technology,
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
func (s *ServerImplementation) GetQuestionById(ctx echo.Context, id int64) error {
	// Получаем вопрос из БД
	question, err := s.questionRepo.GetQuestionByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, openapi.ErrorResponse{
			Message: "Question not found",
			Code:    strPtr("QUESTION_NOT_FOUND"),
		})
	}

	// Преобразуем в формат OpenAPI
	questionDetail := openapi.QuestionDetail{
		Id:            question.ID,
		Title:         question.Title,
		Content:       question.Content,
		Difficulty:    openapi.QuestionDetailDifficulty(question.Difficulty),
		Technology:    question.Technology,
		Options:       question.Options,
		CorrectAnswer: question.CorrectAnswer,
		Explanation:   &question.Explanation,
	}

	return ctx.JSON(http.StatusOK, questionDetail)
}

// strPtr возвращает указатель на строку
func strPtr(s string) *string {
	return &s
}
