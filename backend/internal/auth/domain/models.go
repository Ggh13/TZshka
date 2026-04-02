package domain

import (
	"context"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID   string
	Role Role
}

type TokenService interface {
	GenerateToken(ctx context.Context, role string) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error)
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
