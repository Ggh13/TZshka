package service

import (
	"context"
	"errors"

	sessionModels "TZshka/internal/session/models"
	sessionRepo "TZshka/internal/session/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Service struct {
	repo *sessionRepo.Repository
	log  *zap.Logger
}

func New(repo *sessionRepo.Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) Create(ctx context.Context, name string, creatorID uuid.UUID) (*sessionModels.Session, error) {
	s.log.Info("creating session", zap.String("name", name), zap.String("creatorId", creatorID.String()))

	session := &sessionModels.Session{
		ID:        uuid.New(),
		Name:      name,
		CreatorID: creatorID,
	}

	if err := s.repo.Create(ctx, session); err != nil {
		s.log.Error("failed to create session", zap.Error(err))
		return nil, err
	}

	s.log.Info("session created", zap.String("sessionId", session.ID.String()))
	return session, nil
}

func (s *Service) GetByUser(ctx context.Context, userID uuid.UUID) ([]sessionModels.Session, error) {
	s.log.Info("getting sessions", zap.String("userId", userID.String()))
	return s.repo.GetByCreatorID(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, sessionID uuid.UUID) (*sessionModels.Session, error) {
	return s.repo.GetByID(ctx, sessionID)
}

func (s *Service) Delete(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	s.log.Info("deleting session", zap.String("sessionId", sessionID.String()), zap.String("userId", userID.String()))

	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		s.log.Error("session not found", zap.Error(err))
		return ErrSessionNotFound
	}

	if session.CreatorID != userID {
		s.log.Error("unauthorized to delete session")
		return errors.New("unauthorized")
	}

	return s.repo.Delete(ctx, sessionID)
}
