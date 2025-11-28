package services

import (
	"errors"
	"it_rabotyagi/internal/business/models"
	"net"
	"net/mail"
	"net/url"
	"strings"
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

// isValidEmail проверяет корректность email адреса
func isValidEmail(email string) bool {
	// Используем стандартную библиотеку для парсинга
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Проверяем, что адрес совпадает с оригиналом
	if addr.Address != email {
		return false
	}

	// Разбиваем на локальную часть и домен
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]

	// Отклоняем localhost и локальные домены
	if strings.HasSuffix(domain, ".local") || domain == "localhost" {
		return false
	}

	// Проверяем consecutive dots
	if strings.Contains(email, "..") {
		return false
	}

	return true
}

// isValidURL проверяет корректность URL аватара без whitelist доменов
func isValidURL(urlStr string) bool {
	// Парсим URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Только HTTPS для безопасности
	if u.Scheme != "https" {
		return false
	}

	// Проверяем, что URL не пустой
	if u.Host == "" {
		return false
	}

	// Проверяем, что нет credentials в URL (защита от утечки данных)
	if u.User != nil {
		return false
	}

	hostname := u.Hostname()

	// Блокируем localhost и локальные адреса (защита от SSRF)
	if hostname == "localhost" || strings.HasSuffix(hostname, ".local") {
		return false
	}

	// Проверяем IP адрес - блокируем приватные и loopback адреса
	if ip := net.ParseIP(hostname); ip != nil {
		// Блокируем loopback (127.0.0.0/8, ::1)
		if ip.IsLoopback() {
			return false
		}
		// Блокируем приватные сети (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, fc00::/7)
		if ip.IsPrivate() {
			return false
		}
		// Блокируем link-local адреса (169.254.0.0/16, fe80::/10)
		if ip.IsLinkLocalUnicast() {
			return false
		}
	}

	// Проверяем длину URL (защита от DoS)
	if len(urlStr) > 2048 {
		return false
	}

	return true
}
