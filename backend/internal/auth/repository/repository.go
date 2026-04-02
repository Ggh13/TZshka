package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	authModels "TZshka/internal/auth/models"
)

type Repository interface {
	Create(ctx context.Context, user *authModels.User, passwordHash string) error
	GetByEmail(ctx context.Context, email string) (*authModels.User, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*authModels.User, error)
}

type repo struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func New(db *pgxpool.Pool, log *zap.Logger) Repository {
	return &repo{db: db, log: log}
}

func (r *repo) Create(ctx context.Context, user *authModels.User, passwordHash string) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (email) DO NOTHING
		RETURNING id`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, user.ID, user.Email, passwordHash, user.Role).Scan(&id)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("user already exists")
	}
	if err != nil {
		r.log.Error("failed to create user", zap.Error(err))
		return fmt.Errorf("authRepo.Create: failed to create user: %w", err)
	}

	user.ID = id
	r.log.Info("user created", zap.String("email", user.Email), zap.String("role", user.Role))
	return nil
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*authModels.User, string, error) {
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`

	var user authModels.User
	var passwordHash string
	var createdAt interface{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&passwordHash,
		&user.Role,
		&createdAt,
	)
	if err == pgx.ErrNoRows {
		return nil, "", fmt.Errorf("user not found")
	}
	if err != nil {
		r.log.Error("failed to get user by email", zap.Error(err))
		return nil, "", fmt.Errorf("authRepo.GetByEmail: failed to get user: %w", err)
	}

	return &user, passwordHash, nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*authModels.User, error) {
	query := `SELECT id, email, role, created_at FROM users WHERE id = $1`

	var user authModels.User
	var createdAt interface{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&createdAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		r.log.Error("failed to get user by id", zap.Error(err))
		return nil, fmt.Errorf("authRepo.GetByID: failed to get user: %w", err)
	}

	return &user, nil
}
