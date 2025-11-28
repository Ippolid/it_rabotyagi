package services

import (
	"errors"
	"it_rabotyagi/internal/business/models"
	"regexp"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// ValidateProfileUpdate валидирует данные обновления профиля
func (s *UserService) ValidateProfileUpdate(update *models.UserProfileUpdate) error {
	if update.Username != nil {
		if len(*update.Username) == 0 || len(*update.Username) > 50 {
			return errors.New("username must be 1-50 characters")
		}
	}

	if update.Email != nil {
		if !isValidEmail(*update.Email) {
			return errors.New("invalid email format")
		}
		if len(*update.Email) > 250 {
			return errors.New("email too long")
		}
	}

	if update.Name != nil && len(*update.Name) > 50 {
		return errors.New("name too long")
	}

	if update.Description != nil && len(*update.Description) > 150 {
		return errors.New("description too long")
	}

	return nil
}

// ValidatePassword валидирует новый пароль
func (s *UserService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

// ValidateAvatarURL валидирует URL аватара
func (s *UserService) ValidateAvatarURL(url string) error {
	if !isValidURL(url) {
		return errors.New("invalid avatar URL format")
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidURL(urlStr string) bool {
	urlRegex := regexp.MustCompile(`^https?://[^\s]+$`)
	return urlRegex.MatchString(urlStr)
}
