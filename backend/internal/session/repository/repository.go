package repository

import (
	"context"
	"time"

	sessionModels "TZshka/internal/session/models"
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

func (r *Repository) Create(ctx context.Context, session *sessionModels.Session) error {
	query := `
		INSERT INTO sessions (id, name, creator_id, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, session.ID, session.Name, session.CreatorID, time.Now())
	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*sessionModels.Session, error) {
	query := `
		SELECT id, name, creator_id, created_at
		FROM sessions
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)

	var session sessionModels.Session
	err := row.Scan(&session.ID, &session.Name, &session.CreatorID, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) GetByCreatorID(ctx context.Context, creatorID uuid.UUID) ([]sessionModels.Session, error) {
	query := `
		SELECT id, name, creator_id, created_at
		FROM sessions
		WHERE creator_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []sessionModels.Session
	for rows.Next() {
		var s sessionModels.Session
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatorID, &s.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
