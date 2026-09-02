package usecase

import (
	"context"
	"errors"
	"fmt"

	profileDomain "kopelko-dating-app-backend/internal/modules/profile/domain"
	subscriptionDomain "kopelko-dating-app-backend/internal/modules/subscription/domain"
	"kopelko-dating-app-backend/internal/modules/swipe/domain"
)

type swipeUsecase struct {
	swipeRepo       domain.SwipeRepository
	subscriptionSvc subscriptionDomain.SubscriptionService
	profileRepo     profileDomain.ProfileRepository
	limitSwipe      int
}

func NewSwipeUsecase(
	swipeRepo domain.SwipeRepository,
	subSvc subscriptionDomain.SubscriptionService,
	profileRepo profileDomain.ProfileRepository,
	limitSwipe int,
) domain.SwipeService {
	return &swipeUsecase{
		swipeRepo:       swipeRepo,
		subscriptionSvc: subSvc,
		profileRepo:     profileRepo,
		limitSwipe:      limitSwipe,
	}
}

func (u *swipeUsecase) SwipeUser(ctx context.Context, userID, targetUserID uint, swipeType domain.SwipeType) (*domain.SwipeResponse, error) {
	if userID == targetUserID {
		return nil, errors.New("cannot swipe yourself")
	}

	// Check if already swiped
	alreadySwiped, err := u.swipeRepo.HasSwiped(ctx, userID, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing swipe: %w", err)
	}
	if alreadySwiped {
		return nil, errors.New("you have already swiped on this profile")
	}

	// Check quota
	hasNoQuota, err := u.subscriptionSvc.HasActiveFeature(ctx, userID, "no_swipe_quota")
	if err != nil {
		return nil, fmt.Errorf("failed to verify quota feature: %w", err)
	}

	if !hasNoQuota {
		count, err := u.swipeRepo.GetDailySwipeCount(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to count daily swipes: %w", err)
		}
		if count >= u.limitSwipe {
			return nil, errors.New("daily swipe limit reached. Upgrade to premium for unlimited swipes")
		}
	}

	// Create swipe record
	swipe := &domain.Swipe{
		UserID:       userID,
		TargetUserID: targetUserID,
		SwipeType:    swipeType,
	}
	if err := u.swipeRepo.CreateSwipe(ctx, swipe); err != nil {
		return nil, fmt.Errorf("failed to save swipe: %w", err)
	}

	// Record profile view
	_ = u.profileRepo.RecordView(ctx, userID, targetUserID, &swipe.ID)

	resp := &domain.SwipeResponse{
		Swipe:   swipe,
		IsMatch: false,
	}

	// Check if mutual like (Match)
	if swipeType == domain.SwipeLike {
		targetSwipe, err := u.swipeRepo.GetSwipe(ctx, targetUserID, userID)
		if err == nil && targetSwipe != nil && targetSwipe.SwipeType == domain.SwipeLike {
			match, err := u.swipeRepo.CreateMatch(ctx, userID, targetUserID)
			if err == nil {
				resp.IsMatch = true
				resp.Match = match
			}
		}
	}

	return resp, nil
}
