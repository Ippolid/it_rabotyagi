package jwt

import (
	"errors"
	"fmt"
	"itpath/internal/data/entities"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Константы для типов токенов и времени их жизни
const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour // 30 дней

	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// TokenManager управляет созданием и валидацией JWT токенов.
type TokenManager struct {
	secretKey string
}

// Claims определяет структуру данных, хранимых в JWT токене.
// ИЗМЕНЕНО: Удалены json теги, чтобы jwt.ParseWithClaims работал корректно,
// используя имена полей структуры.
type Claims struct {
	jwt.RegisteredClaims
	TokenType    string
	UserID       int64
	TelegramID   int64
	Role         string
	Subscription string
}

// UserTokenData содержит информацию о пользователе для генерации токена.
type UserTokenData struct {
	UserID       int64
	TelegramID   int64
	Role         entities.UserRole
	Subscription entities.SubscriptionType
}

// NewTokenManager создает новый экземпляр TokenManager.
func NewTokenManager(secretKey string) (*TokenManager, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("secret key cannot be empty")
	}
	return &TokenManager{
		secretKey: secretKey,
	}, nil
}

// GenerateTokens создает пару access и refresh токенов для пользователя.
func (tm *TokenManager) GenerateTokens(data UserTokenData) (accessToken, refreshToken string, err error) {
	// Access token (короткий срок жизни)
	accessClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", data.UserID),
		},
		UserID:       data.UserID,
		TokenType:    tokenTypeAccess,
		TelegramID:   data.TelegramID,
		Role:         string(data.Role),
		Subscription: string(data.Subscription),
	}

	accessToken, err = tm.generateToken(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh token (длительный срок жизни, содержит только ID и TelegramID)
	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", data.UserID),
		},
		UserID:     data.UserID,
		TelegramID: data.TelegramID, // Добавлено для удобства при обновлении
		TokenType:  tokenTypeRefresh,
	}

	refreshToken, err = tm.generateToken(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateToken проверяет подпись токена и возвращает содержащиеся в нем claims.
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
	})

	// Если ошибка есть, но это не ошибка истечения срока, то это серьезная проблема.
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	// Если токен валиден (включая проверку срока) или ошибка - это только истечение срока,
	// мы все равно возвращаем claims. Вызывающий код решит, что делать.
	if token.Valid || errors.Is(err, jwt.ErrTokenExpired) {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// generateToken создает и подписывает JWT токен с заданными claims.
func (tm *TokenManager) generateToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.secretKey))
}
