package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"itpath/internal/business/models"
	"itpath/internal/data/entities"
	"itpath/internal/data/repositories"
	"itpath/internal/logger"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

// Claims представляет JWT claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// ============================================================================
// JWT методы
// ============================================================================

// GenerateToken генерирует JWT токен для пользователя
func (s *AuthService) GenerateToken(user *models.User, provider string) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour) // 7 дней

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	claims := &Claims{
		UserID:   user.ID,
		Email:    email,
		Name:     user.Name,
		Role:     string(user.Role),
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "itpath",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет JWT токен и возвращает claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateCSRFToken генерирует CSRF токен
func (s *AuthService) GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ============================================================================
// Основные методы для работы с пользователями
// ============================================================================

// GetUserByID получает пользователя по ID из базы данных
func (s *AuthService) GetUserByID(id int64) (*models.User, error) {
	userEntity, err := s.userRepo.FindUserByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return models.ConvertToModel(userEntity), nil
}

// UpdateUser обновляет данные пользователя в базе данных
func (s *AuthService) UpdateUser(user *models.User) error {
	userEntity := models.ConvertToEntity(user)
	return s.userRepo.UpdateUser(userEntity)
}

// ============================================================================
// OAuth методы - GitHub
// ============================================================================

// GetOrCreateUserFromGitHub получает или создает пользователя на основе данных GitHub
func (s *AuthService) GetOrCreateUserFromGitHub(githubUser *GitHubUser) (*models.User, error) {
	githubID := strconv.FormatInt(githubUser.ID, 10)

	// Пытаемся найти существующего пользователя по GitHub ID
	userEntity, err := s.userRepo.FindUserByGitHubID(githubID)
	if err == nil {
		// Пользователь найден - обновляем его данные
		s.updateUserFromGitHub(userEntity, githubUser)
		if updateErr := s.userRepo.UpdateUser(userEntity); updateErr != nil {
			logger.Error("Failed to update user from GitHub", zap.Error(updateErr))
		}

		logger.Info("User authenticated via GitHub", zap.Int64("user_id", userEntity.ID))
		return models.ConvertToModel(userEntity), nil
	}

	// Если есть email, проверяем существует ли пользователь с таким email
	if githubUser.Email != nil && *githubUser.Email != "" {
		userEntity, err = s.userRepo.FindUserByEmail(*githubUser.Email)
		if err == nil {
			// Пользователь с таким email уже существует - линкуем GitHub аккаунт
			logger.Info("Linking GitHub account to existing user by email",
				zap.Int64("user_id", userEntity.ID),
				zap.String("email", *githubUser.Email),
				zap.String("github_login", githubUser.Login))

			// Обновляем GitHub ID и другие данные
			userEntity.GitHubID = &githubID
			s.updateUserFromGitHub(userEntity, githubUser)

			if updateErr := s.userRepo.UpdateUser(userEntity); updateErr != nil {
				logger.Error("Failed to link GitHub account", zap.Error(updateErr))
				return nil, fmt.Errorf("failed to link GitHub account: %w", updateErr)
			}

			logger.Info("GitHub account linked to existing user", zap.Int64("user_id", userEntity.ID))
			return models.ConvertToModel(userEntity), nil
		}
	}

	// Пользователь не найден - создаем нового
	logger.Info("Creating new user from GitHub", zap.String("login", githubUser.Login))

	userEntity = &entities.UserEntity{
		GitHubID: &githubID,
		Username: &githubUser.Login,
		Name:     *githubUser.Name,
		Role:     entities.RoleUser,
	}

	s.updateUserFromGitHub(userEntity, githubUser)

	if err := s.userRepo.CreateUser(userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("New user created from GitHub", zap.Int64("user_id", userEntity.ID))
	return models.ConvertToModel(userEntity), nil
}

// updateUserFromGitHub обновляет entity из данных GitHub (только пустые поля)
func (s *AuthService) updateUserFromGitHub(user *entities.UserEntity, githubUser *GitHubUser) {
	// Обновляем name только если оно пустое
	if user.Name == "" {
		if githubUser.Name != nil && *githubUser.Name != "" {
			user.Name = *githubUser.Name
		} else {
			user.Name = githubUser.Login
		}
	}

	// Email - обновляем только если пустой
	if user.Email == nil || *user.Email == "" {
		if githubUser.Email != nil && *githubUser.Email != "" {
			user.Email = githubUser.Email
		}
	}

	// Username - обновляем только если пустой
	if user.Username == nil || *user.Username == "" {
		user.Username = &githubUser.Login
	}

	// Avatar - обновляем только если пустой
	if user.AvatarURL == nil || *user.AvatarURL == "" {
		if githubUser.AvatarURL != "" {
			user.AvatarURL = &githubUser.AvatarURL
		}
	}

	// Description - обновляем только если пустой
	if user.Description == nil || *user.Description == "" {
		var descParts []string

		if githubUser.Bio != nil && *githubUser.Bio != "" {
			descParts = append(descParts, *githubUser.Bio)
		}

		//if githubUser.Company != nil && *githubUser.Company != "" {
		//	descParts = append(descParts, "🏢 "+*githubUser.Company)
		//}
		//
		//if githubUser.Location != nil && *githubUser.Location != "" {
		//	descParts = append(descParts, "📍 "+*githubUser.Location)
		//}
		//
		//if githubUser.Blog != nil && *githubUser.Blog != "" {
		//	descParts = append(descParts, "🔗 "+*githubUser.Blog)
		//}

		if len(descParts) > 0 {
			desc := ""
			for i, part := range descParts {
				if i > 0 {
					desc += " | "
				}
				desc += part
			}
			user.Description = &desc
		}
	}

	logger.Debug("Updated user from GitHub",
		zap.String("login", githubUser.Login),
		zap.Any("name", githubUser.Name),
		zap.Any("email", githubUser.Email),
		zap.Int("public_repos", githubUser.PublicRepos),
		zap.Int("followers", githubUser.Followers))
}

// LinkGitHubAccount связывает GitHub аккаунт с существующим пользователем
func (s *AuthService) LinkGitHubAccount(userID int64, githubUser *GitHubUser) error {
	userEntity, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	githubID := githubUser.Login

	// Проверяем, не занят ли уже этот GitHub ID другим пользователем
	existingUser, err := s.userRepo.FindUserByGitHubID(githubID)
	if err == nil && existingUser.ID != userID {
		return fmt.Errorf("GitHub account already linked to another user")
	}

	// Обновляем данные из GitHub
	userEntity.GitHubID = &githubID
	s.updateUserFromGitHub(userEntity, githubUser)

	logger.Info("GitHub account linked", zap.Int64("user_id", userID), zap.String("github_id", githubID))
	return s.userRepo.UpdateUser(userEntity)
}

// UnlinkGitHubAccount отвязывает GitHub аккаунт от пользователя
func (s *AuthService) UnlinkGitHubAccount(userID int64) error {
	userEntity, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Очищаем GitHub ID
	userEntity.GitHubID = nil

	logger.Info("GitHub account unlinked", zap.Int64("user_id", userID))
	return s.userRepo.UpdateUser(userEntity)
}

// ============================================================================
// OAuth методы - Google
// ============================================================================

// GetOrCreateUserFromGoogle получает или создает пользователя на основе данных Google
func (s *AuthService) GetOrCreateUserFromGoogle(googleUser *GoogleUser) (*models.User, error) {
	googleID := googleUser.ID

	// Пытаемся найти существующего пользователя по Google ID
	userEntity, err := s.userRepo.FindUserByGoogleID(googleID)
	if err == nil {
		// Пользователь найден - обновляем его данные
		s.updateUserFromGoogle(userEntity, googleUser)
		if updateErr := s.userRepo.UpdateUser(userEntity); updateErr != nil {
			logger.Error("Failed to update user from Google", zap.Error(updateErr))
		}

		logger.Info("User authenticated via Google", zap.Int64("user_id", userEntity.ID))
		return models.ConvertToModel(userEntity), nil
	}

	// Проверяем существует ли пользователь с таким email
	if googleUser.Email != "" {
		userEntity, err = s.userRepo.FindUserByEmail(googleUser.Email)
		if err == nil {
			// Пользователь с таким email уже существует - линкуем Google аккаунт
			logger.Info("Linking Google account to existing user by email",
				zap.Int64("user_id", userEntity.ID),
				zap.String("email", googleUser.Email),
				zap.String("google_id", googleUser.ID))

			// Обновляем Google ID и другие данные
			userEntity.GoogleID = &googleID
			s.updateUserFromGoogle(userEntity, googleUser)

			if updateErr := s.userRepo.UpdateUser(userEntity); updateErr != nil {
				logger.Error("Failed to link Google account", zap.Error(updateErr))
				return nil, fmt.Errorf("failed to link Google account: %w", updateErr)
			}

			logger.Info("Google account linked to existing user", zap.Int64("user_id", userEntity.ID))
			return models.ConvertToModel(userEntity), nil
		}
	}

	// Пользователь не найден - создаем нового
	logger.Info("Creating new user from Google", zap.String("email", googleUser.Email))

	userEntity = &entities.UserEntity{
		GoogleID: &googleID,
		Email:    &googleUser.Email,
		Name:     googleUser.Name,
		Role:     entities.RoleUser,
	}

	s.updateUserFromGoogle(userEntity, googleUser)

	if err := s.userRepo.CreateUser(userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("New user created from Google", zap.Int64("user_id", userEntity.ID))
	return models.ConvertToModel(userEntity), nil
}

// updateUserFromGoogle обновляет entity из данных Google (только пустые поля)
func (s *AuthService) updateUserFromGoogle(user *entities.UserEntity, googleUser *GoogleUser) {
	// Обновляем name только если оно пустое
	if user.Name == "" {
		user.Name = googleUser.Name
	}

	// Email - обновляем только если пустой
	if user.Email == nil || *user.Email == "" {
		user.Email = &googleUser.Email
	}

	// Avatar - обновляем только если пустой
	if user.AvatarURL == nil || *user.AvatarURL == "" {
		if googleUser.Picture != "" {
			user.AvatarURL = &googleUser.Picture
		}
	}

	logger.Debug("Updated user from Google",
		zap.String("email", googleUser.Email),
		zap.String("name", googleUser.Name),
		zap.Bool("verified_email", googleUser.VerifiedEmail))
}

// LinkGoogleAccount связывает Google аккаунт с существующим пользователем
func (s *AuthService) LinkGoogleAccount(userID int64, googleUser *GoogleUser) error {
	userEntity, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	googleID := googleUser.ID

	// Проверяем, не занят ли уже этот Google ID другим пользователем
	existingUser, err := s.userRepo.FindUserByGoogleID(googleID)
	if err == nil && existingUser.ID != userID {
		return fmt.Errorf("Google account already linked to another user")
	}

	// Обновляем данные из Google
	userEntity.GoogleID = &googleID
	s.updateUserFromGoogle(userEntity, googleUser)

	logger.Info("Google account linked", zap.Int64("user_id", userID), zap.String("google_id", googleID))
	return s.userRepo.UpdateUser(userEntity)
}

// UnlinkGoogleAccount отвязывает Google аккаунт от пользователя
func (s *AuthService) UnlinkGoogleAccount(userID int64) error {
	userEntity, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Очищаем Google ID
	userEntity.GoogleID = nil

	logger.Info("Google account unlinked", zap.Int64("user_id", userID))
	return s.userRepo.UpdateUser(userEntity)
}

// ============================================================================
// Методы для Telegram (сохранены для совместимости)
// ============================================================================

// LinkTelegramAccount связывает Telegram аккаунт с существующим пользователем
func (s *AuthService) LinkTelegramAccount(userID int64, telegramID string, telegramData map[string]interface{}) error {
	userEntity, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Проверяем, не занят ли уже этот Telegram ID другим пользователем
	existingUser, err := s.userRepo.FindUserByTelegramID(telegramID)
	if err == nil && existingUser.ID != userID {
		return fmt.Errorf("Telegram account already linked to another user")
	}

	// Обновляем Telegram ID
	userEntity.TelegramID = &telegramID

	// Обновляем дополнительные данные из telegramData
	if username, ok := telegramData["username"].(string); ok && username != "" {
		userEntity.Username = &username
	}
	if avatarURL, ok := telegramData["photo_url"].(string); ok && avatarURL != "" {
		userEntity.AvatarURL = &avatarURL
	}
	if firstName, ok := telegramData["first_name"].(string); ok {
		lastName, _ := telegramData["last_name"].(string)
		fullName := firstName
		if lastName != "" {
			fullName += " " + lastName
		}
		userEntity.Name = fullName
	}

	logger.Info("Telegram account linked", zap.Int64("user_id", userID), zap.String("telegram_id", telegramID))
	return s.userRepo.UpdateUser(userEntity)
}

// UnlinkTelegramAccount отвязывает Telegram аккаунт от пользователя
func (s *AuthService) UnlinkTelegramAccount(userID int64) error {
	userEntity, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Очищаем Telegram ID
	userEntity.TelegramID = nil

	logger.Info("Telegram account unlinked", zap.Int64("user_id", userID))
	return s.userRepo.UpdateUser(userEntity)
}
