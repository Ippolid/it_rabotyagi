package services

import (
	"context"
	"fmt"
	"itpath/internal/business"
	"itpath/internal/business/models"
	"itpath/internal/data"
	"itpath/internal/data/entities"
	"itpath/internal/pkg/jwt"
)

type authService struct {
	userRepo        data.UserRepository
	telegramService business.TelegramService
	jwtManager      *jwt.TokenManager
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s *authService) Logout(ctx context.Context, userID int64) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthService(
	userRepo data.UserRepository,
	telegramService business.TelegramService,
	jwtManager *jwt.TokenManager,
) business.AuthService {
	return &authService{
		userRepo:        userRepo,
		telegramService: telegramService,
		jwtManager:      jwtManager,
	}
}

func (s *authService) AuthenticateWithTelegram(ctx context.Context, data models.TelegramAuthData) (*models.AuthResult, error) {
	// 1. Валидируем данные от Telegram
	if err := s.telegramService.ValidateAuthData(data); err != nil {
		return nil, fmt.Errorf("telegram validation failed: %w", err)
	}

	// 2. Ищем существующего пользователя
	userEntity, err := s.userRepo.GetByTelegramID(ctx, data.ID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 3. Создаем или обновляем пользователя
	if userEntity == nil {
		// Создаем нового пользователя
		createReq := entities.CreateUserRequest{
			TelegramID:   data.ID,
			Username:     data.Username,
			FirstName:    data.FirstName,
			LastName:     data.LastName,
			PhotoURL:     data.PhotoURL,
			Role:         entities.RoleUser,
			Subscription: entities.SubscriptionTrial,
		}

		userEntity, err = s.userRepo.Create(ctx, createReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Обновляем данные существующего пользователя
		userEntity, err = s.userRepo.UpdateTelegramData(ctx, data.ID, data.Username, data.FirstName, data.LastName, data.PhotoURL)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	// 4. Конвертируем в бизнес-модель
	user := models.FromEntity(userEntity)

	// 5. Генерируем токены
	accessToken, refreshToken, err := s.jwtManager.GenerateTokens(
		user.ID,
		user.TelegramID,
		getStringValue(user.Username),
		user.FirstName,
		string(user.Role),
		string(user.Subscription),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToPublic(),
		ExpiresIn:    15 * 60, // 15 минут
	}, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if userEntity == nil {
		return nil, fmt.Errorf("user not found")
	}

	return models.FromEntity(userEntity), nil
}

func (s *authService) UpdateProfile(ctx context.Context, userID int64, req models.UpdateProfileRequest) (*models.User, error) {
	// Проверяем существование пользователя
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Валидируем данные
	if err := s.validateUpdateRequest(ctx, userID, req); err != nil {
		return nil, err
	}

	// Конвертируем в запрос для репозитория
	updateReq := entities.UpdateUserRequest{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	// Обновляем пользователя
	userEntity, err := s.userRepo.Update(ctx, userID, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return models.FromEntity(userEntity), nil
}

func (s *authService) ValidateAccess(ctx context.Context, userID int64, requiredRole string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return false, nil
	}

	// Менторы имеют доступ ко всем функциям пользователей
	if requiredRole == "user" {
		return true, nil
	}

	return string(user.Role) == requiredRole, nil
}

func (s *authService) ValidateSubscription(ctx context.Context, userID int64, requiredSubscription string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return false, nil
	}

	// Pro подписка включает все возможности trial
	if requiredSubscription == "trial" {
		return true, nil
	}

	return string(user.Subscription) == requiredSubscription, nil
}

func (s *authService) validateUpdateRequest(ctx context.Context, userID int64, req models.UpdateProfileRequest) error {
	// Проверка email на уникальность
	if req.Email != nil && *req.Email != "" {
		exists, err := s.userRepo.EmailExists(ctx, *req.Email, userID)
		if err != nil {
			return fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			return fmt.Errorf("email already exists")
		}
	}

	// Проверка username на уникальность
	if req.Username != nil && *req.Username != "" {
		exists, err := s.userRepo.UsernameExists(ctx, *req.Username, userID)
		if err != nil {
			return fmt.Errorf("failed to check username: %w", err)
		}
		if exists {
			return fmt.Errorf("username already exists")
		}
	}

	return nil
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
