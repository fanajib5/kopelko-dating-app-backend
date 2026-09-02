package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kopelko-dating-app-backend/internal/modules/swipe/domain"
	"kopelko-dating-app-backend/internal/platform/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type swipePostgresRepo struct {
	db *pgxpool.Pool
}

func NewSwipeRepository(db *pgxpool.Pool) domain.SwipeRepository {
	return &swipePostgresRepo{db: db}
}

func (r *swipePostgresRepo) CreateSwipe(ctx context.Context, s *domain.Swipe) error {
	return r.CreateSwipeWithTx(ctx, r.db, s)
}

func (r *swipePostgresRepo) CreateSwipeWithTx(ctx context.Context, tx database.DBTX, s *domain.Swipe) error {
	query := `
		INSERT INTO swipes (user_id, target_user_id, swipe_type, swipe_date, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_DATE, $4, $5)
		RETURNING id, swipe_date
	`
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	err := tx.QueryRow(ctx, query, s.UserID, s.TargetUserID, s.SwipeType, s.CreatedAt, s.UpdatedAt).Scan(&s.ID, &s.SwipeDate)
	if err != nil {
		return fmt.Errorf("failed to insert swipe: %w", err)
	}
	return nil
}

func (r *swipePostgresRepo) GetDailySwipeCount(ctx context.Context, userID uint) (int, error) {
	return r.GetDailySwipeCountWithTx(ctx, r.db, userID)
}

func (r *swipePostgresRepo) GetDailySwipeCountWithTx(ctx context.Context, tx database.DBTX, userID uint) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM swipes
		WHERE user_id = $1 AND swipe_date = CURRENT_DATE
	`
	var count int
	err := tx.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count daily swipes: %w", err)
	}
	return count, nil
}

func (r *swipePostgresRepo) HasSwiped(ctx context.Context, userID, targetUserID uint) (bool, error) {
	return r.HasSwipedWithTx(ctx, r.db, userID, targetUserID)
}

func (r *swipePostgresRepo) HasSwipedWithTx(ctx context.Context, tx database.DBTX, userID, targetUserID uint) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM swipes WHERE user_id = $1 AND target_user_id = $2
		)
	`
	var exists bool
	err := tx.QueryRow(ctx, query, userID, targetUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check swipe status: %w", err)
	}
	return exists, nil
}

func (r *swipePostgresRepo) GetSwipe(ctx context.Context, userID, targetUserID uint) (*domain.Swipe, error) {
	return r.GetSwipeWithTx(ctx, r.db, userID, targetUserID)
}

func (r *swipePostgresRepo) GetSwipeWithTx(ctx context.Context, tx database.DBTX, userID, targetUserID uint) (*domain.Swipe, error) {
	query := `
		SELECT id, user_id, target_user_id, swipe_type, swipe_date, created_at, updated_at, deleted_at
		FROM swipes
		WHERE user_id = $1 AND target_user_id = $2 AND deleted_at IS NULL
		ORDER BY id DESC
		LIMIT 1
	`
	var s domain.Swipe
	err := tx.QueryRow(ctx, query, userID, targetUserID).Scan(
		&s.ID, &s.UserID, &s.TargetUserID, &s.SwipeType, &s.SwipeDate, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get swipe: %w", err)
	}
	return &s, nil
}

func (r *swipePostgresRepo) CreateMatch(ctx context.Context, user1ID, user2ID uint) (*domain.Match, error) {
	return r.CreateMatchWithTx(ctx, r.db, user1ID, user2ID)
}

func (r *swipePostgresRepo) CreateMatchWithTx(ctx context.Context, tx database.DBTX, user1ID, user2ID uint) (*domain.Match, error) {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
		INSERT INTO matches (user1_id, user2_id, matched_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user1_id, user2_id) DO UPDATE SET matched_at = EXCLUDED.matched_at
		RETURNING id, user1_id, user2_id, matched_at
	`
	var m domain.Match
	err := tx.QueryRow(ctx, query, user1ID, user2ID, time.Now()).Scan(&m.ID, &m.User1ID, &m.User2ID, &m.MatchedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create match: %w", err)
	}
	return &m, nil
}

func (r *swipePostgresRepo) GetMatch(ctx context.Context, user1ID, user2ID uint) (*domain.Match, error) {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
		SELECT id, user1_id, user2_id, matched_at
		FROM matches
		WHERE user1_id = $1 AND user2_id = $2
		LIMIT 1
	`
	var m domain.Match
	err := r.db.QueryRow(ctx, query, user1ID, user2ID).Scan(&m.ID, &m.User1ID, &m.User2ID, &m.MatchedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query match: %w", err)
	}
	return &m, nil
}
