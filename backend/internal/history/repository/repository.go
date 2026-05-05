package repository

import (
	"context"
<<<<<<< HEAD
	"time"

	historyModels "TZshka/internal/history/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
=======
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	historyModels "TZshka/internal/history/models"
)

type Repository interface {
	Create(ctx context.Context, correction *historyModels.Correction) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]historyModels.Correction, error)
}

type repo struct {
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
	db  *pgxpool.Pool
	log *zap.Logger
}

<<<<<<< HEAD
func New(db *pgxpool.Pool, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) Save(ctx context.Context, sessionID uuid.UUID, inputContent string, responseData map[string]interface{}) error {
	query := `
		INSERT INTO corrections_history (id, session_id, input_content, response_data, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), sessionID, inputContent, responseData, time.Now())
	return err
}

func (r *Repository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]historyModels.Correction, error) {
	query := `
		SELECT id, session_id, input_content, response_data, created_at
		FROM corrections_history
		WHERE session_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
=======
func New(db *pgxpool.Pool, log *zap.Logger) Repository {
	return &repo{db: db, log: log}
}

func (r *repo) Create(ctx context.Context, correction *historyModels.Correction) error {
	respDataJSON, err := json.Marshal(correction.ResponseData)
	if err != nil {
		r.log.Error("failed to marshal response data", zap.Error(err))
		return fmt.Errorf("failed to marshal response data: %w", err)
	}

	query := `
		INSERT INTO corrections_history (id, session_id, input_content, response_data, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id`

	var id uuid.UUID
	err = r.db.QueryRow(ctx, query, correction.ID, correction.SessionID, correction.InputContent, respDataJSON).Scan(&id)
	if err != nil {
		r.log.Error("failed to create correction", zap.Error(err))
		return fmt.Errorf("historyRepo.Create: failed to create correction: %w", err)
	}

	correction.ID = id
	r.log.Info("correction created", zap.String("id", correction.ID.String()), zap.String("sessionId", correction.SessionID.String()))
	return nil
}

func (r *repo) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]historyModels.Correction, error) {
	query := `SELECT id, session_id, input_content, response_data, created_at FROM corrections_history WHERE session_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		r.log.Error("failed to get corrections", zap.Error(err))
		return nil, fmt.Errorf("historyRepo.GetBySessionID: failed to get corrections: %w", err)
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
	}
	defer rows.Close()

	var corrections []historyModels.Correction
	for rows.Next() {
<<<<<<< HEAD
		var c historyModels.Correction
		if err := rows.Scan(&c.ID, &c.SessionID, &c.InputContent, &c.ResponseData, &c.CreatedAt); err != nil {
			return nil, err
		}
		corrections = append(corrections, c)
	}
	return corrections, nil
}
=======
		var correction historyModels.Correction
		var responseDataJSON []byte
		var createdAt interface{}

		if err := rows.Scan(&correction.ID, &correction.SessionID, &correction.InputContent, &responseDataJSON, &createdAt); err != nil {
			r.log.Error("failed to scan correction", zap.Error(err))
			return nil, fmt.Errorf("failed to scan correction: %w", err)
		}

		if err := json.Unmarshal(responseDataJSON, &correction.ResponseData); err != nil {
			r.log.Error("failed to unmarshal response data", zap.Error(err))
			return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
		}

		corrections = append(corrections, correction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if corrections == nil {
		return []historyModels.Correction{}, nil
	}

	return corrections, nil
}
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
