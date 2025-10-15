package services

import (
	"fmt"
	"itpath/internal/business/models"
	"itpath/internal/data/entities"
	"itpath/internal/data/repositories"
	"log"
	"strconv"

	"github.com/go-pkgz/auth/token"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// GetOrCreateUser получает или создает пользователя на основе данных от провайдера
func (s *AuthService) GetOrCreateUser(claims token.User) (*models.User, error) {
	var userEntity *entities.UserEntity
	var err error

	// Определяем провайдера и ищем пользователя
	switch {
	case claims.ID != "" && len(claims.ID) > 0:
		// Попробуем найти по разным провайдерам
		if userEntity, err = s.findUserByProvider(claims); err == nil {

			fmt.Pr
			return s.convertToModel(userEntity), nil
		}

		// Если не нашли - создаем нового
		userEntity = s.createUserEntityFromClaims(claims)
		if err := s.userRepo.Create(userEntity); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		return s.convertToModel(userEntity), nil
	}

	return nil, fmt.Errorf("invalid user claims")
}

// findUserByProvider ищет пользователя по ID от провайдера
func (s *AuthService) findUserByProvider(claims token.User) (*entities.UserEntity, error) {
	// Telegram
	if claims.ID != "" {
		user, err := s.userRepo.FindByTelegramID(claims.ID)
		if err == nil {
			return user, nil
		}
	}

	// Google (если ID начинается с google_ или содержит email)
	if claims.Email != "" {
		user, err := s.userRepo.FindByGoogleID(claims.ID)
		if err == nil {
			return user, nil
		}
	}

	// GitHub
	user, err := s.userRepo.FindByGitHubID(claims.ID)
	if err == nil {
		return user, nil
	}

	return nil, fmt.Errorf("user not found")
}

// createUserEntityFromClaims создает entity пользователя из claims
func (s *AuthService) createUserEntityFromClaims(claims token.User) *entities.UserEntity {
	user := &entities.UserEntity{
		Name: claims.Name,
		Role: entities.RoleUser,
	}

	// Определяем провайдера и заполняем соответствующие поля
	if claims.Email != "" {
		user.Email = &claims.Email
	}

	// Telegram
	if claims.ID != "" && claims.Attributes != nil {
		if _, ok := claims.Attributes["telegram"]; ok {
			user.TelegramID = &claims.ID
			if username, ok := claims.Attributes["username"].(string); ok {
				user.Username = &username
			}
		}
	}

	// Google
	if claims.Email != "" && claims.Attributes != nil {
		if _, ok := claims.Attributes["google"]; ok {
			user.GoogleID = &claims.ID
		}
	}

	// GitHub
	if claims.Attributes != nil {
		if _, ok := claims.Attributes["github"]; ok {
			user.GitHubID = &claims.ID
			if username, ok := claims.Attributes["username"].(string); ok {
				user.Username = &username
			}
		}
	}

	if claims.Picture != "" {
		user.AvatarURL = &claims.Picture
	}

	return user
}

// GetUserByID получает пользователя по ID
func (s *AuthService) GetUserByID(id int64) (*models.User, error) {
	userEntity, err := s.userRepo.FindByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.convertToModel(userEntity), nil
}

// convertToModel конвертирует entity в бизнес-модель
func (s *AuthService) convertToModel(userEntity *entities.UserEntity) *models.User {
	user := &models.User{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Role:      string(userEntity.Role),
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	if userEntity.TelegramID != nil {
		user.TelegramID = userEntity.TelegramID
	}
	if userEntity.GoogleID != nil {
		user.GoogleID = userEntity.GoogleID
	}
	if userEntity.GitHubID != nil {
		user.GitHubID = userEntity.GitHubID
	}
	if userEntity.Email != nil {
		user.Email = userEntity.Email
	}
	if userEntity.Username != nil {
		user.Username = userEntity.Username
	}
	if userEntity.AvatarURL != nil {
		user.AvatarURL = userEntity.AvatarURL
	}
	if userEntity.Description != nil {
		user.Description = userEntity.Description
	}
	if userEntity.SubscriptionType != nil {
		subType := string(*userEntity.SubscriptionType)
		user.SubscriptionType = &subType
	}
	if userEntity.SubscriptionExpiresAt != nil {
		user.SubscriptionExpiresAt = userEntity.SubscriptionExpiresAt
	}

	return user
}

// UpdateUser обновляет данные пользователя
func (s *AuthService) UpdateUser(user *models.User) error {
	userEntity := &entities.UserEntity{
		ID:          user.ID,
		TelegramID:  user.TelegramID,
		GoogleID:    user.GoogleID,
		GitHubID:    user.GitHubID,
		Email:       user.Email,
		Username:    user.Username,
		Name:        user.Name,
		AvatarURL:   user.AvatarURL,
		Description: user.Description,
		Role:        entities.UserRole(user.Role),
	}

	if user.SubscriptionType != nil {
		subType := entities.SubscriptionType(*user.SubscriptionType)
		userEntity.SubscriptionType = &subType
	}

	if user.SubscriptionExpiresAt != nil {
		userEntity.SubscriptionExpiresAt = user.SubscriptionExpiresAt
	}

	return s.userRepo.Update(userEntity)
}

// ClaimsUpdater обновляет claims после успешной OAuth авторизации
// Этот метод вызывается go-pkgz/auth и позволяет сохранить пользователя в БД
func (s *AuthService) ClaimsUpdater(claims token.Claims) token.Claims {
	// Пропускаем handshake токены (промежуточные токены OAuth процесса)
	if claims.Handshake != nil {
		log.Printf("[AUTH] Skipping handshake token")
		return claims
	}

	// Проверяем, что User не nil (должен быть заполнен после успешной OAuth авторизации)
	if claims.User == nil {
		log.Printf("[AUTH ERROR] claims.User is nil (not a handshake token)")
		return claims
	}

	//log.Printf("[AUTH] Processing claims for user: %s (ID: %s)", claims.User.Name, claims.User.ID)

	// Получаем или создаем пользователя в БД
	fmt.Println(claims.User, "        dfdsfsdfsfsfs")
	user, err := s.GetOrCreateUserFromClaims(*claims.User)
	if err != nil {
		log.Printf("[AUTH ERROR] Failed to get/create user: %v", err)
		return claims
	}

	log.Printf("[AUTH] User saved to DB with ID: %d", user.ID)

	// Добавляем ID из БД в claims для дальнейшего использования
	if claims.User.Attributes == nil {
		claims.User.Attributes = make(map[string]interface{})
	}
	claims.User.Attributes["db_user_id"] = user.ID

	return claims
}

// GetOrCreateUserFromClaims получает или создает пользователя на основе OAuth claims
func (s *AuthService) GetOrCreateUserFromClaims(claims token.User) (*models.User, error) {
	// Определяем провайдера по audience
	provider := s.detectProvider(claims)
	log.Printf("[AUTH] Detected provider: %s for user: %s", provider, claims.Name)
	log.Printf("[AUTH] Claims ID: %s", claims.ID)
	log.Printf("[AUTH] Claims Attributes: %+v", claims.Attributes)

	// Извлекаем ID пользователя без префикса провайдера
	userID := s.extractProviderUserID(claims.ID, provider)
	log.Printf("[AUTH] Extracted user ID: %s", userID)

	var userEntity *entities.UserEntity
	var err error

	// Для GitHub ищем по login из attributes, а не по хешу
	if provider == "github" && claims.Attributes != nil {
		if login, ok := claims.Attributes["login"].(string); ok && login != "" {
			log.Printf("[AUTH] Searching GitHub user by login: %s", login)
			userEntity, err = s.userRepo.FindByGitHubID(login)
			if err == nil && userEntity != nil {
				log.Printf("[AUTH] User found in DB: ID=%d", userEntity.ID)
				// Обновляем данные пользователя (например, avatar, name)
				s.updateUserFromClaims(userEntity, claims, provider)
				if updateErr := s.userRepo.Update(userEntity); updateErr != nil {
					log.Printf("[AUTH] Failed to update user: %v", updateErr)
				}
				return s.convertToModel(userEntity), nil
			}
		}
	}

	// Пытаемся найти пользователя по ID провайдера для других провайдеров
	if provider != "github" {
		switch provider {
		case "telegram":
			userEntity, err = s.userRepo.FindByTelegramID(userID)
		case "google":
			userEntity, err = s.userRepo.FindByGoogleID(userID)
		default:
			// Пытаемся найти по email если есть
			if claims.Email != "" {
				userEntity, err = s.findByEmail(claims.Email)
			}
		}

		// Если пользователь найден - обновляем данные и возвращаем
		if err == nil && userEntity != nil {
			log.Printf("[AUTH] User found in DB: ID=%d", userEntity.ID)
			// Обновляем данные пользователя (например, avatar, name)
			s.updateUserFromClaims(userEntity, claims, provider)
			if updateErr := s.userRepo.Update(userEntity); updateErr != nil {
				log.Printf("[AUTH] Failed to update user: %v", updateErr)
			}
			return s.convertToModel(userEntity), nil
		}
	}

	// Пользователь не найден - создаем нового
	log.Printf("[AUTH] Creating new user for %s", claims.Name)
	userEntity = s.createUserEntityFromOAuthClaims(claims, provider)
	log.Printf("[AUTH] New user entity to create: GitHubID=%v, Username=%v, Email=%v",
		userEntity.GitHubID, userEntity.Username, userEntity.Email)

	if err := s.userRepo.Create(userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("[AUTH] New user created with ID: %d", userEntity.ID)
	return s.convertToModel(userEntity), nil
}

// detectProvider определяет OAuth провайдера по claims
func (s *AuthService) detectProvider(claims token.User) string {
	// Проверяем по ID (go-pkgz/auth добавляет префикс провайдера к ID)
	if len(claims.ID) > 7 && claims.ID[:7] == "github_" {
		return "github"
	}
	if len(claims.ID) > 7 && claims.ID[:7] == "google_" {
		return "google"
	}
	if len(claims.ID) > 9 && claims.ID[:9] == "telegram_" {
		return "telegram"
	}

	// Проверяем audience
	if claims.Audience != "" {
		return claims.Audience
	}

	// Проверяем attributes
	if claims.Attributes != nil {
		if provider, ok := claims.Attributes["provider"].(string); ok {
			return provider
		}
	}

	// Определяем по email (Google обычно возвращает email)
	if claims.Email != "" {
		return "google"
	}

	return "unknown"
}

// createUserEntityFromOAuthClaims создает entity из OAuth claims
func (s *AuthService) createUserEntityFromOAuthClaims(claims token.User, provider string) *entities.UserEntity {
	user := &entities.UserEntity{
		Name: claims.Name,
		Role: entities.RoleUser,
	}

	// Извлекаем правильный ID пользователя из claims.ID (убираем префикс провайдера)
	userID := s.extractProviderUserID(claims.ID, provider)

	// Заполняем поля в зависимости от провайдера
	switch provider {
	case "telegram":
		user.TelegramID = &userID
		if claims.Attributes != nil {
			if username, ok := claims.Attributes["username"].(string); ok && username != "" {
				user.Username = &username
			}
		}
	case "google":
		user.GoogleID = &userID
		if claims.Email != "" {
			user.Email = &claims.Email
		}
	case "github":
		// Для GitHub пытаемся получить login (username) из attributes
		if claims.Attributes != nil {
			if login, ok := claims.Attributes["login"].(string); ok && login != "" {
				user.GitHubID = &login
				user.Username = &login
			} else {
				// Если login нет, используем ID
				user.GitHubID = &userID
			}
			if email, ok := claims.Attributes["email"].(string); ok && email != "" {
				user.Email = &email
			}
			if username, ok := claims.Attributes["username"].(string); ok && username != "" {
				user.Username = &username
			}
		} else {
			user.GitHubID = &userID
		}
	}

	if claims.Picture != "" {
		user.AvatarURL = &claims.Picture
	}

	return user
}

// extractProviderUserID извлекает ID пользователя, убирая префикс провайдера
func (s *AuthService) extractProviderUserID(fullID string, provider string) string {
	prefix := provider + "_"
	if len(fullID) > len(prefix) && fullID[:len(prefix)] == prefix {
		return fullID[len(prefix):]
	}
	return fullID
}

// updateUserFromClaims обновляет entity из OAuth claims
func (s *AuthService) updateUserFromClaims(user *entities.UserEntity, claims token.User, provider string) {
	// Обновляем имя если изменилось
	if claims.Name != "" && claims.Name != user.Name {
		user.Name = claims.Name
	}

	// Обновляем специфичные данные провайдера
	switch provider {
	case "telegram":
		if user.TelegramID == nil {
			user.TelegramID = &claims.ID
		}
		if claims.Attributes != nil {
			if username, ok := claims.Attributes["username"].(string); ok && username != "" {
				user.Username = &username
			}
		}
		// Обновляем avatar если изменился
		if claims.Picture != "" {
			user.AvatarURL = &claims.Picture
		}
	case "google":
		if user.GoogleID == nil {
			user.GoogleID = &claims.ID
		}
		// Обновляем email если есть
		if claims.Email != "" && (user.Email == nil || *user.Email != claims.Email) {
			user.Email = &claims.Email
		}
		// Обновляем avatar если изменился
		if claims.Picture != "" {
			user.AvatarURL = &claims.Picture
		}
	case "github":
		// Обновляем GitHubID по login
		if claims.Attributes != nil {
			if login, ok := claims.Attributes["login"].(string); ok && login != "" {
				if user.GitHubID == nil || *user.GitHubID != login {
					user.GitHubID = &login
					user.Username = &login
				}
			}

			// Обновляем bio (description)
			if bio, ok := claims.Attributes["bio"].(string); ok && bio != "" {
				user.Description = &bio
			}

			// Обновляем email если есть
			if email, ok := claims.Attributes["email"].(string); ok && email != "" {
				user.Email = &email
			}

			// Обновляем avatar_url из GitHub attributes
			if avatarURL, ok := claims.Attributes["avatar_url"].(string); ok && avatarURL != "" {
				user.AvatarURL = &avatarURL
			}
		}

		// Если avatar_url не был в attributes, используем claims.Picture
		if user.AvatarURL == nil && claims.Picture != "" {
			user.AvatarURL = &claims.Picture
		}
	}
}

// findByEmail ищет пользователя по email
func (s *AuthService) findByEmail(email string) (*entities.UserEntity, error) {
	return s.userRepo.FindByEmail(email)
}

// GetUserFromToken получает пользователя из token.User (для использования в handlers)
func (s *AuthService) GetUserFromToken(tokenUser token.User) (*models.User, error) {
	// Пытаемся получить db_user_id из attributes
	if tokenUser.Attributes != nil {
		if dbUserID, ok := tokenUser.Attributes["db_user_id"]; ok {
			var userID int64
			switch v := dbUserID.(type) {
			case int64:
				userID = v
			case float64:
				userID = int64(v)
			case string:
				id, err := strconv.ParseInt(v, 10, 64)
				if err == nil {
					userID = id
				}
			}

			if userID > 0 {
				return s.GetUserByID(userID)
			}
		}
	}

	// Если нет db_user_id, пытаемся найти по провайдеру
	provider := s.detectProvider(tokenUser)
	var userEntity *entities.UserEntity
	var err error

	switch provider {
	case "telegram":
		userEntity, err = s.userRepo.FindByTelegramID(tokenUser.ID)
	case "google":
		userEntity, err = s.userRepo.FindByGoogleID(tokenUser.ID)
	case "github":
		userEntity, err = s.userRepo.FindByGitHubID(tokenUser.ID)
	default:
		if tokenUser.Email != "" {
			userEntity, err = s.userRepo.FindByEmail(tokenUser.Email)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.convertToModel(userEntity), nil
}
