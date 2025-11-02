package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	secret          string
	tokenDuration   time.Duration
	refreshDuration time.Duration
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(secret string, tokenDuration, refreshDuration int) *AuthService {
	return &AuthService{
		secret:          secret,
		tokenDuration:   time.Duration(tokenDuration) * time.Minute,
		refreshDuration: time.Duration(refreshDuration) * time.Minute,
	}
}

// GenerateTokens генерирует access и refresh токены
func (s *AuthService) GenerateTokens(userID int, nickname, role string) (accessToken, refreshToken string, expiresIn int, err error) {
	// Access token
	accessClaims := Claims{
		UserID:   userID,
		Nickname: nickname,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", 0, err
	}

	// Refresh token
	refreshClaims := Claims{
		UserID:   userID,
		Nickname: nickname,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, int(s.tokenDuration.Seconds()), nil
}

// ValidateToken проверяет токен и возвращает claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshAccessToken обновляет access токен по refresh токену
func (s *AuthService) RefreshAccessToken(refreshToken string) (accessToken string, expiresIn int, err error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", 0, err
	}

	// Генерируем новый access токен
	accessClaims := Claims{
		UserID:   claims.UserID,
		Nickname: claims.Nickname,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(s.secret))
	if err != nil {
		return "", 0, err
	}

	return accessToken, int(s.tokenDuration.Seconds()), nil
}

// GetRefreshTokenExpiration возвращает время истечения refresh токена
func (s *AuthService) GetRefreshTokenExpiration() time.Duration {
	return s.refreshDuration
}

// GetAccessTokenExpiration возвращает время истечения access токена
func (s *AuthService) GetAccessTokenExpiration() time.Duration {
	return s.tokenDuration
}
