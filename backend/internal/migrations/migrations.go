package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Name string
	Up   func(ctx context.Context, tx *pgxpool.Pool) error
}

var migrations = []Migration{
	{
		Name: "001_create_users",
		Up: func(ctx context.Context, tx *pgxpool.Pool) error {
			query := `
			CREATE TABLE IF NOT EXISTS users (
				id UUID PRIMARY KEY,
				login VARCHAR(50) UNIQUE NOT NULL,
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT NOW()
			);`
			_, err := tx.Exec(ctx, query)
			return err
		},
	},
	{
		Name: "002_create_sessions",
		Up: func(ctx context.Context, tx *pgxpool.Pool) error {
			query := `
			CREATE TABLE IF NOT EXISTS sessions (
				id UUID PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				creator_id UUID NOT NULL,
				created_at TIMESTAMP DEFAULT NOW()
			);`
			_, err := tx.Exec(ctx, query)
			return err
		},
	},
	{
		Name: "003_create_corrections_history",
		Up: func(ctx context.Context, tx *pgxpool.Pool) error {
			query := `
			CREATE TABLE IF NOT EXISTS corrections_history (
				id UUID PRIMARY KEY,
				session_id UUID NOT NULL,
				input_content TEXT NOT NULL,
				response_data JSONB NOT NULL,
				created_at TIMESTAMP DEFAULT NOW()
			);`
			_, err := tx.Exec(ctx, query)
			return err
		},
	},
}

func Run(ctx context.Context, pool *pgxpool.Pool) error {
	for _, m := range migrations {
		if err := m.Up(ctx, pool); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.Name, err)
		}
	}
	return nil
}
