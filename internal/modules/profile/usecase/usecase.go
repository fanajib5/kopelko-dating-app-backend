package usecase

import (
	"context"
	"errors"
	"fmt"

	"kopelko-dating-app-backend/internal/modules/profile/domain"
	subscriptionDomain "kopelko-dating-app-backend/internal/modules/subscription/domain"
)

type profileUsecase struct {
	repo            domain.ProfileRepository
	subscriptionSvc subscriptionDomain.SubscriptionService
	limitSwipe      int
}

func NewProfileUsecase(
	repo domain.ProfileRepository,
	subSvc subscriptionDomain.SubscriptionService,
	limitSwipe int,
) domain.ProfileService {
	return &profileUsecase{
		repo:            repo,
		subscriptionSvc: subSvc,
		limitSwipe:      limitSwipe,
	}
}

func (u *profileUsecase) GetMyProfile(ctx context.Context, userID uint) (*domain.Profile, error) {
	return u.repo.GetByUserID(ctx, userID)
}

func (u *profileUsecase) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	return u.repo.Create(ctx, profile)
}

func (u *profileUsecase) GetRandomProfiles(ctx context.Context, currentUserID uint) ([]domain.Profile, error) {
	// Check if user is unlimited
	hasNoQuota, err := u.subscriptionSvc.HasActiveFeature(ctx, currentUserID, "no_swipe_quota")
	if err != nil {
		return nil, fmt.Errorf("could not verify user subscription: %w", err)
	}

	limit := u.limitSwipe
	if !hasNoQuota {
		count, err := u.repo.GetDailyViewCount(ctx, currentUserID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch daily view count: %w", err)
		}
		if count >= u.limitSwipe {
			return nil, errors.New("daily profile view limit reached. Upgrade to premium for unlimited profiles")
		}
		limit = u.limitSwipe - count
	}

	profiles, err := u.repo.GetRandomProfiles(ctx, currentUserID, limit)
	if err != nil {
		return nil, err
	}

	// Record view for fetched profiles
	for _, p := range profiles {
		_ = u.repo.RecordView(ctx, currentUserID, p.UserID, nil)
	}

	return profiles, nil
}
