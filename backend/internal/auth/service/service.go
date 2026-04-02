package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	authModels "TZshka/internal/auth/models"
	authRepo "TZshka/internal/auth/repository"
	"TZshka/pkg/registr"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRole  = errors.New("invalid role: must be 'admin' or 'user'")
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
)

type Service struct {
	repo         authRepo.Repository
	tokenService *registr.TokenService
	log          *zap.Logger
}

func New(repo authRepo.Repository, secret string, log *zap.Logger) *Service {
	return &Service{
		repo:         repo,
		tokenService: registr.NewTokenService(secret),
		log:          log,
	}
}

func (s *Service) Register(ctx context.Context, email, password, role string) (*authModels.User, error) {
	s.log.Info("registering user", zap.String("email", email), zap.String("role", role))

	if role != "admin" && role != "user" {
		s.log.Error("invalid role", zap.String("role", role))
		return nil, ErrInvalidRole
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("authService.Register: failed to hash password: %w", err)
	}

	user := &authModels.User{
		ID:    uuid.New(),
		Email: email,
		Role:  role,
	}

	if err := s.repo.Create(ctx, user, string(hash)); err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		if strings.Contains(err.Error(), "user already exists") || strings.Contains(err.Error(), "duplicate") {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("authService.Register: failed to create user: %w", err)
	}

	s.log.Info("user registered successfully", zap.String("userId", user.ID.String()))
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*authModels.User, string, error) {
	s.log.Info("logging in user", zap.String("email", email))

	user, passwordHash, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("user not found", zap.Error(err))
		return nil, "", ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		s.log.Error("invalid password", zap.Error(err))
		return nil, "", ErrInvalidCreds
	}

	token, err := s.tokenService.GenerateToken(ctx, user.Role)
	if err != nil {
		s.log.Error("failed to generate token", zap.Error(err))
		return nil, "", fmt.Errorf("authService.Login: failed to sign token: %w", err)
	}

	s.log.Info("user logged in successfully", zap.String("userId", user.ID.String()))
	return user, token, nil
}
