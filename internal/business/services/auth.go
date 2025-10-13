package services

import (
	"context"
	"errors"
	"fmt"
	"itpath/internal/business"
	"itpath/internal/business/models"
	"itpath/internal/data"
	"itpath/internal/data/entities"
	"itpath/internal/pkg/jwt"
	"log"
)

type authService struct {
	userRepo        data.UserRepository
	telegramService business.TelegramService
	jwtManager      *jwt.TokenManager
}

// NewAuthService создает новый экземпляр сервиса аутентификации.
func NewAuthService(userRepo data.UserRepository, telegramService business.TelegramService, jwtManager *jwt.TokenManager) business.AuthService {
	return &authService{
		userRepo:        userRepo,
		telegramService: telegramService,
		jwtManager:      jwtManager,
	}
}

func (a *authService) AuthenticateWithTelegram(ctx context.Context, data models.TelegramAuthData) (*models.AuthResult, error) {
	// 1. Валидация данных от Telegram
	if err := a.telegramService.ValidateAuthData(data); err != nil {
		return nil, fmt.Errorf("telegram data validation failed: %w", err)
	}

	// 2. Поиск пользователя или создание нового
	user, err := a.userRepo.GetByTelegramID(ctx, data.ID)
	if err != nil {
		// Если пользователь не найден, создаем его
		if errors.Is(err, entities.ErrUserNotFound) {
			createReq := entities.CreateUserRequest{
				TelegramID: data.ID,
				FirstName:  data.FirstName,
				LastName:   &data.LastName,
				Username:   &data.Username,
				PhotoURL:   &data.PhotoURL,
			}
			user, err = a.userRepo.CreateUser(ctx, createReq)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get user by telegram id: %w", err)
		}
	}

	// 3. Генерация токенов
	tokenData := jwt.UserTokenData{
		UserID:       user.ID,
		TelegramID:   user.TelegramID,
		Role:         user.Role,
		Subscription: user.Subscription,
	}

	accessToken, refreshToken, err := a.jwtManager.GenerateTokens(tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResult, error) {
	// 1. Валидация refresh токена
	claims, err := a.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type, expected refresh")
	}

	// 2. Получение пользователя из БД
	log.Printf("claims: %+v", claims)
	user, err := a.userRepo.GetByTelegramID(ctx, claims.TelegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for token refresh: %w", err)
	}

	// 3. Генерация новой пары токенов
	tokenData := jwt.UserTokenData{
		UserID:       user.ID,
		TelegramID:   user.TelegramID,
		Role:         user.Role,
		Subscription: user.Subscription,
	}

	newAccessToken, newRefreshToken, err := a.jwtManager.GenerateTokens(tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return &models.AuthResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *authService) Logout(ctx context.Context, userID int64) error {
	// Для JWT-based аутентификации, выход обычно обрабатывается на клиенте
	// путем удаления токенов. Если требуется принудительный отзыв,
	// необходимо реализовать механизм черного списка токенов (например, в Redis).
	return nil
}

func (a *authService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	userEntity, err := a.userRepo.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	// Конвертация из entity в business model
	// Конвертация из entity в business model
	user := models.FromEntity(userEntity)
	return user, nil
}

func (a *authService) UpdateProfile(ctx context.Context, userID int64, req models.UpdateProfileRequest) (*models.User, error) {
	updateReq := entities.UpdateUserRequest{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		PhotoURL:  req.PhotoURL,
	}

	updatedUserEntity, err := a.userRepo.Update(ctx, userID, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	user := models.FromEntity(updatedUserEntity)

	return user, nil
}
