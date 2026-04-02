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
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				role VARCHAR(50) NOT NULL,
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
