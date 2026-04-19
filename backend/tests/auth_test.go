package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authModels "TZshka/internal/auth/models"
	"github.com/gin-gonic/gin"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/register", func(c *gin.Context) {
		var req authModels.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "invalid request"},
			})
			return
		}

		if req.Login == "" {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "login is required"},
			})
			return
		}

		if req.Email == "" {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "email is required"},
			})
			return
		}

		if req.Password == "" {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "password is required"},
			})
			return
		}

		c.JSON(http.StatusCreated, authModels.UserResponse{
			User: authModels.User{
				Login: req.Login,
				Email: req.Email,
			},
		})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid register",
			body: map[string]string{
				"login":    "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing login",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			body: map[string]string{
				"login":    "testuser",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			body: map[string]string{
				"login": "testuser",
				"email": "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/login", func(c *gin.Context) {
		var req authModels.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "invalid request"},
			})
			return
		}

		if req.Login == "" {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "login is required"},
			})
			return
		}

		if req.Password == "" {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "password is required"},
			})
			return
		}

		if req.Login == "testuser" && req.Password == "password123" {
			c.JSON(http.StatusOK, authModels.TokenResponse{
				Token: "valid-token",
			})
			return
		}

		c.JSON(http.StatusUnauthorized, authModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "UNAUTHORIZED", Message: "invalid credentials"},
		})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid login",
			body: map[string]string{
				"login":    "testuser",
				"password": "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing login",
			body: map[string]string{
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			body: map[string]string{
				"login": "testuser",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			body: map[string]string{
				"login":    "wronguser",
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
