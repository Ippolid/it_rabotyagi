package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"itpath/internal/logger"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuthService управляет OAuth аутентификацией
type OAuthService struct {
	githubConfig *oauth2.Config
	googleConfig *oauth2.Config
	states       map[string]*OAuthState
	statesMu     sync.RWMutex
}

// OAuthState хранит информацию о состоянии OAuth запроса
type OAuthState struct {
	Provider  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// GitHubUser представляет данные пользователя GitHub
type GitHubUser struct {
	Login             string  `json:"login"`
	ID                int64   `json:"id"`
	NodeID            string  `json:"node_id"`
	AvatarURL         string  `json:"avatar_url"`
	GravatarID        string  `json:"gravatar_id"`
	URL               string  `json:"url"`
	HTMLURL           string  `json:"html_url"`
	FollowersURL      string  `json:"followers_url"`
	FollowingURL      string  `json:"following_url"`
	GistsURL          string  `json:"gists_url"`
	StarredURL        string  `json:"starred_url"`
	SubscriptionsURL  string  `json:"subscriptions_url"`
	OrganizationsURL  string  `json:"organizations_url"`
	ReposURL          string  `json:"repos_url"`
	EventsURL         string  `json:"events_url"`
	ReceivedEventsURL string  `json:"received_events_url"`
	Type              string  `json:"type"`
	SiteAdmin         bool    `json:"site_admin"`
	Name              *string `json:"name"`
	Company           *string `json:"company"`
	Blog              *string `json:"blog"`
	Location          *string `json:"location"`
	Email             *string `json:"email"`
	Hireable          *bool   `json:"hireable"`
	Bio               *string `json:"bio"`
	TwitterUsername   *string `json:"twitter_username"`
	PublicRepos       int     `json:"public_repos"`
	PublicGists       int     `json:"public_gists"`
	Followers         int     `json:"followers"`
	Following         int     `json:"following"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// GitHubEmail представляет email пользователя GitHub
type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

// GoogleUser представляет данные пользователя Google
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// NewOAuthService создает новый OAuthService
func NewOAuthService(githubClientID, githubClientSecret, googleClientID, googleClientSecret, redirectURL string) *OAuthService {
	service := &OAuthService{
		states: make(map[string]*OAuthState),
	}

	// Настройка GitHub OAuth
	if githubClientID != "" && githubClientSecret != "" {
		service.githubConfig = &oauth2.Config{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURL:  redirectURL + "/auth/github/callback",
			Scopes:       []string{"user:email", "read:user"}, // Запрашиваем доступ к email и профилю
			Endpoint:     github.Endpoint,
		}
	}

	// Настройка Google OAuth
	if googleClientID != "" && googleClientSecret != "" {
		service.googleConfig = &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  redirectURL + "/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	}

	// Запускаем очистку устаревших состояний
	go service.cleanupExpiredStates()

	return service
}

// GetGitHubAuthURL возвращает URL для авторизации через GitHub
func (s *OAuthService) GetGitHubAuthURL() (string, error) {
	if s.githubConfig == nil {
		return "", fmt.Errorf("github oauth not configured")
	}

	state, err := s.generateState("github")
	if err != nil {
		return "", err
	}

	return s.githubConfig.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

// GetGoogleAuthURL возвращает URL для авторизации через Google
func (s *OAuthService) GetGoogleAuthURL() (string, error) {
	if s.googleConfig == nil {
		return "", fmt.Errorf("google oauth not configured")
	}

	state, err := s.generateState("google")
	if err != nil {
		return "", err
	}

	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

// HandleGitHubCallback обрабатывает callback от GitHub
func (s *OAuthService) HandleGitHubCallback(ctx context.Context, state, code string) (*GitHubUser, error) {
	if s.githubConfig == nil {
		return nil, fmt.Errorf("github oauth not configured")
	}

	// Проверяем state
	if !s.validateState(state, "github") {
		return nil, fmt.Errorf("invalid state parameter")
	}

	// Обмениваем code на токен
	token, err := s.githubConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Получаем информацию о пользователе
	user, err := s.getGitHubUser(ctx, token)

	logger.Debug("HandleGitHubCallback", zap.Any("user", user))

	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Если email не указан в профиле, пытаемся получить из API emails
	if user.Email == nil || *user.Email == "" {
		email, _ := s.getGitHubPrimaryEmail(ctx, token)
		if email != "" {
			user.Email = &email
		}
	}

	return user, nil
}

// HandleGoogleCallback обрабатывает callback от Google
func (s *OAuthService) HandleGoogleCallback(ctx context.Context, state, code string) (*GoogleUser, error) {
	if s.googleConfig == nil {
		return nil, fmt.Errorf("google oauth not configured")
	}

	// Проверяем state
	if !s.validateState(state, "google") {
		return nil, fmt.Errorf("invalid state parameter")
	}

	// Обмениваем code на токен
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Получаем информацию о пользователе
	user, err := s.getGoogleUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}

// getGitHubUser получает информацию о пользователе GitHub
func (s *OAuthService) getGitHubUser(ctx context.Context, token *oauth2.Token) (*GitHubUser, error) {
	client := s.githubConfig.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api error: %s", string(body))
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// getGitHubPrimaryEmail получает основной email пользователя GitHub
func (s *OAuthService) getGitHubPrimaryEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	client := s.githubConfig.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get emails")
	}

	var emails []GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	// Ищем основной verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// Если основной не найден, берем первый verified
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}

// getGoogleUser получает информацию о пользователе Google
func (s *OAuthService) getGoogleUser(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := s.googleConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google api error: %s", string(body))
	}

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// generateState генерирует уникальный state для OAuth запроса
func (s *OAuthService) generateState(provider string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)

	s.statesMu.Lock()
	defer s.statesMu.Unlock()

	s.states[state] = &OAuthState{
		Provider:  provider,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	return state, nil
}

// validateState проверяет валидность state
func (s *OAuthService) validateState(state, provider string) bool {
	s.statesMu.RLock()
	defer s.statesMu.RUnlock()

	oauthState, exists := s.states[state]
	if !exists {
		return false
	}

	if oauthState.Provider != provider {
		return false
	}

	if time.Now().After(oauthState.ExpiresAt) {
		return false
	}

	// Удаляем использованный state
	go func() {
		s.statesMu.Lock()
		defer s.statesMu.Unlock()
		delete(s.states, state)
	}()

	return true
}

// cleanupExpiredStates периодически удаляет устаревшие states
func (s *OAuthService) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.statesMu.Lock()
		now := time.Now()
		for state, oauthState := range s.states {
			if now.After(oauthState.ExpiresAt) {
				delete(s.states, state)
			}
		}
		s.statesMu.Unlock()
	}
}
