package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"itpath/internal/business"
	"itpath/internal/business/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

type telegramService struct {
	botToken string
}

func NewTelegramService(botToken string) business.TelegramService {
	return &telegramService{
		botToken: botToken,
	}
}

func (s *telegramService) ValidateAuthData(data models.TelegramAuthData) error {
	if data.Hash == "" {
		return fmt.Errorf("hash is missing")
	}

	if data.ID == 0 {
		return fmt.Errorf("user ID is missing")
	}

	// Создаем data-check-string
	checkString := s.createDataCheckString(data)

	// Создаем secret key
	secretKey := sha256.Sum256([]byte(s.botToken))

	// Вычисляем HMAC
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(checkString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if data.Hash != expectedHash {
		return fmt.Errorf("invalid hash")
	}

	// Проверяем время (не старше 1 дня)
	if time.Now().Unix()-data.AuthDate > 86400 {
		return fmt.Errorf("authorization data is too old")
	}

	return nil
}

func (s *telegramService) createDataCheckString(data models.TelegramAuthData) string {
	var pairs []string

	pairs = append(pairs, "auth_date="+strconv.FormatInt(data.AuthDate, 10))
	pairs = append(pairs, "first_name="+data.FirstName)
	pairs = append(pairs, "id="+strconv.FormatInt(data.ID, 10))

	if data.LastName != "" {
		pairs = append(pairs, "last_name="+data.LastName)
	}
	if data.PhotoURL != "" {
		pairs = append(pairs, "photo_url="+data.PhotoURL)
	}
	if data.Username != "" {
		pairs = append(pairs, "username="+data.Username)
	}

	sort.Strings(pairs)
	return strings.Join(pairs, "\n")
}
