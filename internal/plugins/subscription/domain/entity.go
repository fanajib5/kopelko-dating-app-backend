package domain

import (
	"context"
	"time"
)

type PremiumFeature struct {
	ID          uint   `json:"id"`
	FeatureName string `json:"feature_name"`
	Description string `json:"description"`
}

type Subscription struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	FeatureID uint      `json:"feature_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	AutoRenew bool      `json:"auto_renew"`
}

type SubscriptionRepository interface {
	GetFeatureByName(ctx context.Context, featureName string) (*PremiumFeature, error)
	GetFeatureByID(ctx context.Context, featureID uint) (*PremiumFeature, error)
	CreateOrRenewSubscription(ctx context.Context, sub *Subscription) error
	HasActiveFeature(ctx context.Context, userID uint, featureName string) (bool, error)
	GetActiveSubscriptions(ctx context.Context, userID uint) ([]Subscription, error)
}

type SubscriptionService interface {
	Subscribe(ctx context.Context, userID uint, featureName string) (*Subscription, error)
	HasActiveFeature(ctx context.Context, userID uint, featureName string) (bool, error)
}
