package dto

import "itpath/internal/business/models"

// TelegramAuthRequest - запрос авторизации через Telegram
type TelegramAuthRequest struct {
	ID        int64  `json:"id" form:"id" binding:"required"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name"`
	Username  string `json:"username" form:"username"`
	PhotoURL  string `json:"photo_url" form:"photo_url"`
	AuthDate  int64  `json:"auth_date" form:"auth_date" binding:"required"`
	Hash      string `json:"hash" form:"hash" binding:"required"`
}

// AuthResponse - ответ авторизации
type AuthResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	User         *models.PublicUser `json:"user"`
	ExpiresIn    int64              `json:"expires_in"`
}

// UpdateProfileRequest - запрос обновления профиля
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=1,max=255"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=255"`
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50,alphanum"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
}

// ToBusinessModel - конвертация в бизнес-модель
func (r *TelegramAuthRequest) ToBusinessModel() models.TelegramAuthData {
	return models.TelegramAuthData{
		ID:        r.ID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Username:  r.Username,
		PhotoURL:  r.PhotoURL,
		AuthDate:  r.AuthDate,
		Hash:      r.Hash,
	}
}

func (r *UpdateProfileRequest) ToBusinessModel() models.UpdateProfileRequest {
	return models.UpdateProfileRequest{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Username:  r.Username,
		Email:     r.Email,
	}
}
