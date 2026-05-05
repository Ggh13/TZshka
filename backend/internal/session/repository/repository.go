package repository

import (
	"context"
<<<<<<< HEAD
	"time"

	sessionModels "TZshka/internal/session/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
=======
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	sessionModels "TZshka/internal/session/models"
)

type Repository interface {
	Create(ctx context.Context, session *sessionModels.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*sessionModels.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]sessionModels.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
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
=======
func New(db *pgxpool.Pool, log *zap.Logger) Repository {
	return &repo{db: db, log: log}
}

func (r *repo) Create(ctx context.Context, session *sessionModels.Session) error {
	query := `
		INSERT INTO sessions (id, name, creator_id, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, session.ID, session.Name, session.CreatorID).Scan(&id)
	if err != nil {
		r.log.Error("failed to create session", zap.Error(err))
		return fmt.Errorf("sessionRepo.Create: failed to create session: %w", err)
	}

	session.ID = id
	r.log.Info("session created", zap.String("sessionId", session.ID.String()), zap.String("name", session.Name))
	return nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*sessionModels.Session, error) {
	query := `SELECT id, name, creator_id, created_at FROM sessions WHERE id = $1`

	var session sessionModels.Session
	var createdAt interface{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.Name,
		&session.CreatorID,
		&createdAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		r.log.Error("failed to get session by id", zap.Error(err))
		return nil, fmt.Errorf("sessionRepo.GetByID: failed to get session: %w", err)
	}

	return &session, nil
}

func (r *repo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]sessionModels.Session, error) {
	query := `SELECT id, name, creator_id, created_at FROM sessions WHERE creator_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.log.Error("failed to get user sessions", zap.Error(err))
		return nil, fmt.Errorf("sessionRepo.GetByUserID: failed to get sessions: %w", err)
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
	}
	defer rows.Close()

	var sessions []sessionModels.Session
	for rows.Next() {
<<<<<<< HEAD
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
=======
		var session sessionModels.Session
		var createdAt interface{}
		if err := rows.Scan(&session.ID, &session.Name, &session.CreatorID, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.log.Error("failed to delete session", zap.Error(err))
		return fmt.Errorf("sessionRepo.Delete: failed to delete session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	r.log.Info("session deleted", zap.String("sessionId", id.String()))
	return nil
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
}
