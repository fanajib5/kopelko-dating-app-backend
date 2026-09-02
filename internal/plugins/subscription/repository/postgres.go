package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kopelko-dating-app-backend/internal/plugins/subscription/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionPostgresRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) domain.SubscriptionRepository {
	return &subscriptionPostgresRepo{db: db}
}

func (r *subscriptionPostgresRepo) GetFeatureByName(ctx context.Context, featureName string) (*domain.PremiumFeature, error) {
	query := `SELECT id, feature_name, COALESCE(description, '') FROM premium_features WHERE feature_name = $1 LIMIT 1`
	var feat domain.PremiumFeature
	err := r.db.QueryRow(ctx, query, featureName).Scan(&feat.ID, &feat.FeatureName, &feat.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("premium feature not found")
		}
		return nil, fmt.Errorf("failed to get feature: %w", err)
	}
	return &feat, nil
}

func (r *subscriptionPostgresRepo) GetFeatureByID(ctx context.Context, featureID uint) (*domain.PremiumFeature, error) {
	query := `SELECT id, feature_name, COALESCE(description, '') FROM premium_features WHERE id = $1 LIMIT 1`
	var feat domain.PremiumFeature
	err := r.db.QueryRow(ctx, query, featureID).Scan(&feat.ID, &feat.FeatureName, &feat.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("premium feature not found")
		}
		return nil, fmt.Errorf("failed to get feature: %w", err)
	}
	return &feat, nil
}

func (r *subscriptionPostgresRepo) CreateOrRenewSubscription(ctx context.Context, sub *domain.Subscription) error {
	// Upsert subscription for user: if already exists, extend or overwrite end_date
	query := `
		INSERT INTO subscriptions (user_id, feature_id, start_date, end_date, auto_renew)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) 
		DO UPDATE SET feature_id = EXCLUDED.feature_id, 
		              start_date = EXCLUDED.start_date, 
		              end_date = EXCLUDED.end_date, 
		              auto_renew = EXCLUDED.auto_renew
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query, sub.UserID, sub.FeatureID, sub.StartDate, sub.EndDate, sub.AutoRenew).Scan(&sub.ID)
	if err != nil {
		return fmt.Errorf("failed to create/renew subscription: %w", err)
	}
	return nil
}

func (r *subscriptionPostgresRepo) HasActiveFeature(ctx context.Context, userID uint, featureName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM subscriptions s
			JOIN premium_features pf ON s.feature_id = pf.id
			WHERE s.user_id = $1 
			  AND pf.feature_name = $2 
			  AND s.end_date > $3
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, featureName, time.Now()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check active feature: %w", err)
	}
	return exists, nil
}

func (r *subscriptionPostgresRepo) GetActiveSubscriptions(ctx context.Context, userID uint) ([]domain.Subscription, error) {
	query := `
		SELECT id, user_id, feature_id, start_date, end_date, auto_renew
		FROM subscriptions
		WHERE user_id = $1 AND end_date > $2
	`
	rows, err := r.db.Query(ctx, query, userID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.FeatureID, &s.StartDate, &s.EndDate, &s.AutoRenew); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}
