package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	sessionModels "TZshka/internal/session/models"
	sessionRepo "TZshka/internal/session/repository"
)

var (
	ErrSessionNotFound = fmt.Errorf("session not found")
	ErrUnauthorized    = fmt.Errorf("unauthorized")
)

type Service struct {
	repo sessionRepo.Repository
	log  *zap.Logger
}

func New(repo sessionRepo.Repository, log *zap.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) CreateSession(ctx context.Context, name string, creatorID uuid.UUID) (*sessionModels.Session, error) {
	s.log.Info("creating session", zap.String("name", name), zap.String("creatorId", creatorID.String()))

	session := &sessionModels.Session{
		ID:        uuid.New(),
		Name:      name,
		CreatorID: creatorID,
	}

	if err := s.repo.Create(ctx, session); err != nil {
		s.log.Error("failed to create session", zap.Error(err))
		return nil, fmt.Errorf("sessionService.Create: failed to create session: %w", err)
	}

	s.log.Info("session created successfully", zap.String("sessionId", session.ID.String()))
	return session, nil
}

func (s *Service) GetSession(ctx context.Context, id uuid.UUID) (*sessionModels.Session, error) {
	s.log.Info("getting session", zap.String("sessionId", id.String()))

	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("session not found", zap.Error(err))
		return nil, ErrSessionNotFound
	}

	return session, nil
}

func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]sessionModels.Session, error) {
	s.log.Info("getting user sessions", zap.String("userId", userID.String()))

	sessions, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user sessions", zap.Error(err))
		return nil, fmt.Errorf("sessionService.GetUserSessions: failed to get sessions: %w", err)
	}

	return sessions, nil
}

func (s *Service) DeleteSession(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	s.log.Info("deleting session", zap.String("sessionId", id.String()), zap.String("userId", userID.String()))

	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrSessionNotFound
	}

	if session.CreatorID != userID {
		s.log.Error("unauthorized to delete session", zap.Error(err))
		return ErrUnauthorized
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete session", zap.Error(err))
		return fmt.Errorf("sessionService.Delete: failed to delete session: %w", err)
	}

	s.log.Info("session deleted successfully", zap.String("sessionId", id.String()))
	return nil
}

func (s *Service) JoinSession(ctx context.Context, id uuid.UUID) (*sessionModels.Session, error) {
	s.log.Info("joining session", zap.String("sessionId", id.String()))

	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("session not found", zap.Error(err))
		return nil, ErrSessionNotFound
	}

	s.log.Info("user joined session", zap.String("sessionId", id.String()))
	return session, nil
}
