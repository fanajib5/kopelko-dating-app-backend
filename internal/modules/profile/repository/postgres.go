package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"kopelko-dating-app-backend/internal/modules/profile/domain"
	"kopelko-dating-app-backend/internal/platform/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type profilePostgresRepo struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) domain.ProfileRepository {
	return &profilePostgresRepo{db: db}
}

func (r *profilePostgresRepo) Create(ctx context.Context, p *domain.Profile) error {
	return r.CreateWithTx(ctx, r.db, p)
}

func (r *profilePostgresRepo) CreateWithTx(ctx context.Context, tx database.DBTX, p *domain.Profile) error {
	query := `
		INSERT INTO profiles (user_id, name, age, bio, gender, location, interests, photos, is_premium, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	err := tx.QueryRow(ctx, query,
		p.UserID, p.Name, p.Age, p.Bio, p.Gender, p.Location, p.Interests, p.Photos, p.IsPremium, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.ID)

	if err != nil {
		return fmt.Errorf("failed to insert profile: %w", err)
	}
	return nil
}

func (r *profilePostgresRepo) GetByUserID(ctx context.Context, userID uint) (*domain.Profile, error) {
	query := `
		SELECT p.id, p.user_id, p.name, p.age, COALESCE(p.bio, ''), p.gender, COALESCE(p.location, ''), 
		       COALESCE(p.interests, '{}'), COALESCE(p.photos, '{}'), p.is_premium, u.is_verified, 
		       p.created_at, p.updated_at, p.deleted_at
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1 AND p.deleted_at IS NULL
		LIMIT 1
	`
	var p domain.Profile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Age, &p.Bio, &p.Gender, &p.Location,
		&p.Interests, &p.Photos, &p.IsPremium, &p.IsVerified, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("profile not found")
		}
		return nil, fmt.Errorf("failed to query profile: %w", err)
	}
	return &p, nil
}

func (r *profilePostgresRepo) Update(ctx context.Context, p *domain.Profile) error {
	query := `
		UPDATE profiles
		SET name = $1, age = $2, bio = $3, gender = $4, location = $5, 
		    interests = $6, photos = $7, is_premium = $8, updated_at = $9
		WHERE user_id = $10 AND deleted_at IS NULL
	`
	p.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		p.Name, p.Age, p.Bio, p.Gender, p.Location, p.Interests, p.Photos, p.IsPremium, p.UpdatedAt, p.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}

func (r *profilePostgresRepo) RecordView(ctx context.Context, userID, viewedUserID uint, swipeID *uint) error {
	return r.RecordViewWithTx(ctx, r.db, userID, viewedUserID, swipeID)
}

func (r *profilePostgresRepo) RecordViewWithTx(ctx context.Context, tx database.DBTX, userID, viewedUserID uint, swipeID *uint) error {
	query := `
		INSERT INTO profile_views (user_id, viewed_user_id, swipe_id, view_date, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_DATE, $4, $4)
		ON CONFLICT (user_id, viewed_user_id, view_date) DO NOTHING
	`
	now := time.Now()
	_, err := tx.Exec(ctx, query, userID, viewedUserID, swipeID, now)
	if err != nil {
		return fmt.Errorf("failed to record profile view: %w", err)
	}
	return nil
}

func (r *profilePostgresRepo) GetDailyViewCount(ctx context.Context, userID uint) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM profile_views
		WHERE user_id = $1 AND view_date = CURRENT_DATE
	`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count daily views: %w", err)
	}
	return count, nil
}

func (r *profilePostgresRepo) GetRandomProfiles(ctx context.Context, currentUserID uint, filter domain.DiscoveryFilter) ([]domain.Profile, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT p.id, p.user_id, p.name, p.age, COALESCE(p.bio, ''), p.gender, COALESCE(p.location, ''), 
		       COALESCE(p.interests, '{}'), COALESCE(p.photos, '{}'), p.is_premium, u.is_verified, 
		       p.created_at, p.updated_at, p.deleted_at
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id != $1
		  AND p.deleted_at IS NULL
		  AND p.user_id NOT IN (
		      SELECT pv.viewed_user_id 
		      FROM profile_views pv 
		      WHERE pv.user_id = $1 AND pv.view_date = CURRENT_DATE
		  )
	`)

	args := []any{currentUserID}
	argIdx := 2

	if filter.Gender != nil && *filter.Gender != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND p.gender = $%d", argIdx))
		args = append(args, *filter.Gender)
		argIdx++
	}

	if filter.MinAge != nil && *filter.MinAge > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" AND p.age >= $%d", argIdx))
		args = append(args, *filter.MinAge)
		argIdx++
	}

	if filter.MaxAge != nil && *filter.MaxAge > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" AND p.age <= $%d", argIdx))
		args = append(args, *filter.MaxAge)
		argIdx++
	}

	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY RANDOM() LIMIT $%d", argIdx))
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	args = append(args, limit)

	rows, err := r.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query random profiles: %w", err)
	}
	defer rows.Close()

	var profiles []domain.Profile
	for rows.Next() {
		var p domain.Profile
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Name, &p.Age, &p.Bio, &p.Gender, &p.Location,
			&p.Interests, &p.Photos, &p.IsPremium, &p.IsVerified, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}
