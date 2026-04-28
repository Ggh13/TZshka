package repository

import (
	"context"
	"time"

	historyModels "TZshka/internal/history/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
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
	}
	defer rows.Close()

	var corrections []historyModels.Correction
	for rows.Next() {
		var c historyModels.Correction
		if err := rows.Scan(&c.ID, &c.SessionID, &c.InputContent, &c.ResponseData, &c.CreatedAt); err != nil {
			return nil, err
		}
		corrections = append(corrections, c)
	}
	return corrections, nil
}
