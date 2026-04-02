package registr

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

var (
	adminUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userUUID  = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey []byte
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secretKey: []byte(secret),
	}
}

func (s *TokenService) GenerateToken(ctx context.Context, role string) (string, error) {
	var userID uuid.UUID
	switch role {
	case RoleAdmin:
		userID = adminUUID
	case RoleUser:
		userID = userUUID
	default:
		return "", ErrInvalidRole
	}

	claims := JWTClaims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *TokenService) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
