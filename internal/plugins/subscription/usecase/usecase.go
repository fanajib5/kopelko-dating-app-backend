package usecase

import (
	"context"
	"fmt"
	"time"

	"kopelko-dating-app-backend/internal/plugins/subscription/domain"
)

type subscriptionUsecase struct {
	repo domain.SubscriptionRepository
}

func NewSubscriptionUsecase(repo domain.SubscriptionRepository) domain.SubscriptionService {
	return &subscriptionUsecase{repo: repo}
}

func (u *subscriptionUsecase) Subscribe(ctx context.Context, userID uint, featureName string) (*domain.Subscription, error) {
	feature, err := u.repo.GetFeatureByName(ctx, featureName)
	if err != nil {
		return nil, fmt.Errorf("invalid feature: %w", err)
	}

	now := time.Now()
	// Default 30 days subscription
	endDate := now.AddDate(0, 1, 0)

	sub := &domain.Subscription{
		UserID:    userID,
		FeatureID: feature.ID,
		StartDate: now,
		EndDate:   endDate,
		AutoRenew: true,
	}

	if err := u.repo.CreateOrRenewSubscription(ctx, sub); err != nil {
		return nil, fmt.Errorf("could not activate subscription: %w", err)
	}

	return sub, nil
}

func (u *subscriptionUsecase) HasActiveFeature(ctx context.Context, userID uint, featureName string) (bool, error) {
	return u.repo.HasActiveFeature(ctx, userID, featureName)
}
