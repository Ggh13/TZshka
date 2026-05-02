package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	historyModels "TZshka/internal/history/models"
	historyRepo "TZshka/internal/history/repository"
)

var (
	ErrCorrectionNotFound = fmt.Errorf("correction not found")
)

type Service struct {
	repo historyRepo.Repository
	log *zap.Logger
}

func New(repo historyRepo.Repository, log *zap.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) CreateCorrection(ctx context.Context, sessionID uuid.UUID, inputContent string, responseData any) (*historyModels.Correction, error) {
	s.log.Info("creating correction", zap.String("sessionId", sessionID.String()))

	correction := &historyModels.Correction{
		ID:           uuid.New(),
		SessionID:    sessionID,
		InputContent: inputContent,
		ResponseData: responseData,
	}

	if err := s.repo.Create(ctx, correction); err != nil {
		s.log.Error("failed to create correction", zap.Error(err))
		return nil, fmt.Errorf("historyService.CreateCorrection: failed to create correction: %w", err)
	}

	s.log.Info("correction created successfully", zap.String("id", correction.ID.String()))
	return correction, nil
}

func (s *Service) GetCorrectionsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]historyModels.Correction, error) {
	s.log.Info("getting corrections for session", zap.String("sessionId", sessionID.String()))

	corrections, err := s.repo.GetBySessionID(ctx, sessionID)
	if err != nil {
		s.log.Error("failed to get corrections", zap.Error(err))
		return nil, fmt.Errorf("historyService.GetCorrectionsBySessionID: failed to get corrections: %w", err)
	}

	s.log.Info("corrections fetched successfully", zap.Int("count", len(corrections)))
	return corrections, nil
}