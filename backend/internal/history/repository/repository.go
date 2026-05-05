package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

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

func (r *Repository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]map[string]interface{}, error) {
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

	var results []map[string]interface{}
	for rows.Next() {
		var id, sessID uuid.UUID
		var inputContent string
		var responseData map[string]interface{}
		var createdAt time.Time

		if err := rows.Scan(&id, &sessID, &inputContent, &responseData, &createdAt); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":           id,
			"sessionId":    sessID,
			"inputContent": inputContent,
			"responseData": responseData,
			"createdAt":    createdAt,
		})
	}
	return results, nil
}
