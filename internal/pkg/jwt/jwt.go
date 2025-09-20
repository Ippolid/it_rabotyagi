package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secretKey string
}

type Claims struct {
	UserID       int64  `json:"user_id"`
	TelegramID   int64  `json:"telegram_id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	Role         string `json:"role"`
	Subscription string `json:"subscription"`
	jwt.RegisteredClaims
	LastName any
}

func NewTokenManager(secretKey string) *TokenManager {
	return &TokenManager{
		secretKey: secretKey,
	}
}

func (tm *TokenManager) GenerateTokens(userID, telegramID int64, username, firstName, role, subscription string) (accessToken, refreshToken string, err error) {
	// Access token (короткий срок жизни)
	accessClaims := Claims{
		UserID:       userID,
		TelegramID:   telegramID,
		Username:     username,
		FirstName:    firstName,
		Role:         role,
		Subscription: subscription,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh token (длительный срок жизни)
	refreshClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)), // 7 дней
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
