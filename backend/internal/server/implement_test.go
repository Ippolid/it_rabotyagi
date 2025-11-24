package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockPermissionRepo struct {
	allowed bool
	err     error
}

func (m mockPermissionRepo) IsEditor(ctx context.Context, userID int) (bool, error) {
	return m.allowed, m.err
}

func TestEnsureEditor(t *testing.T) {
	tests := []struct {
		name           string
		setUser        bool
		userID         int
		role           string
		permissionRepo PermissionChecker
		wantAllowed    bool
		wantStatus     int
	}{
		{
			name:        "no user in context -> unauthorized",
			setUser:     false,
			wantAllowed: false,
			wantStatus:  http.StatusUnauthorized,
		},
		{
			name:           "no role in context -> unauthorized",
			setUser:        true,
			userID:         1,
			permissionRepo: mockPermissionRepo{},
			wantAllowed:    false,
			wantStatus:     http.StatusUnauthorized,
		},
		{
			name:           "admin role bypass",
			setUser:        true,
			userID:         1,
			role:           "admin",
			permissionRepo: mockPermissionRepo{allowed: false},
			wantAllowed:    true,
			wantStatus:     http.StatusOK,
		},
		{
			name:           "editor allowed in table",
			setUser:        true,
			userID:         2,
			role:           "user",
			permissionRepo: mockPermissionRepo{allowed: true},
			wantAllowed:    true,
			wantStatus:     http.StatusOK,
		},
		{
			name:           "editor forbidden",
			setUser:        true,
			userID:         3,
			role:           "user",
			permissionRepo: mockPermissionRepo{allowed: false},
			wantAllowed:    false,
			wantStatus:     http.StatusForbidden,
		},
		{
			name:           "permission check error",
			setUser:        true,
			userID:         4,
			role:           "user",
			permissionRepo: mockPermissionRepo{err: errors.New("db error")},
			wantAllowed:    false,
			wantStatus:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.setUser {
				c.Set(UserIDKey, tt.userID)
				if tt.role != "" {
					c.Set(RoleKey, tt.role)
				}
			}

			impl := &ServerImplementation{
				permissionRepo: tt.permissionRepo,
			}

			allowed := impl.ensureEditor(c)
			if allowed != tt.wantAllowed {
				t.Fatalf("expected allowed=%v, got %v", tt.wantAllowed, allowed)
			}

			status := rec.Code
			if status == 0 {
				status = http.StatusOK
			}
			if status != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, status)
			}
		})
	}
}
