package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kopelko-dating-app-backend/internal/core/identity/domain"
	"kopelko-dating-app-backend/internal/platform/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userPostgresRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userPostgresRepo{db: db}
}

func (r *userPostgresRepo) Create(ctx context.Context, u *domain.User) error {
	return r.CreateWithTx(ctx, r.db, u)
}

func (r *userPostgresRepo) CreateWithTx(ctx context.Context, tx database.DBTX, u *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	err := tx.QueryRow(ctx, query, u.Email, u.PasswordHash, u.IsVerified, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *userPostgresRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, is_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}
	return &u, nil
}

func (r *userPostgresRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, is_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to query user by id: %w", err)
	}
	return &u, nil
}

func (r *userPostgresRepo) SetVerified(ctx context.Context, id uint, isVerified bool) error {
	query := `
		UPDATE users
		SET is_verified = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, isVerified, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update user verification: %w", err)
	}
	return nil
}
