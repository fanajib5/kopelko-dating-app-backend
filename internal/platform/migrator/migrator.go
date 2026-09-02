package migrator

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigration(ctx context.Context, db *pgxpool.Pool, schemaPath string) error {
	slog.Info("Running database migrations...", "schema_path", schemaPath)

	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("failed to execute schema migration: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}

func RunSeeder(ctx context.Context, db *pgxpool.Pool, seederPath string) error {
	slog.Info("Running database seeders...", "seeder_path", seederPath)

	content, err := os.ReadFile(seederPath)
	if err != nil {
		return fmt.Errorf("failed to read seeder file: %w", err)
	}

	_, err = db.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("failed to execute database seeder: %w", err)
	}

	slog.Info("Database seeders completed successfully")
	return nil
}
