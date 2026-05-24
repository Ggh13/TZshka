package service

import (
	"context"

	historyRepo "TZshka/internal/history/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	repo *historyRepo.Repository
	log  *zap.Logger
}

func New(repo *historyRepo.Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) SaveCorrection(ctx context.Context, sessionID, userID uuid.UUID, inputContent string, responseData map[string]interface{}) error {
	s.log.Info("saving correction", zap.String("session_id", sessionID.String()), zap.String("user_id", userID.String()))
	return s.repo.Save(ctx, sessionID, userID, inputContent, responseData)
}

func (s *Service) GetHistory(ctx context.Context, sessionID uuid.UUID) ([]map[string]interface{}, error) {
	s.log.Info("getting history", zap.String("session_id", sessionID.String()))
	return s.repo.GetBySessionID(ctx, sessionID)
}
