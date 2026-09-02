package usecase

import (
	"context"
	"errors"
	"fmt"

	"kopelko-dating-app-backend/internal/core/hook"
	"kopelko-dating-app-backend/internal/plugins/profile/domain"
)

type profileUsecase struct {
	repo       domain.ProfileRepository
	hooks      hook.HookManager
	limitSwipe int
}

func NewProfileUsecase(
	repo domain.ProfileRepository,
	hooks hook.HookManager,
	limitSwipe int,
) domain.ProfileService {
	return &profileUsecase{
		repo:       repo,
		hooks:      hooks,
		limitSwipe: limitSwipe,
	}
}

func (u *profileUsecase) GetMyProfile(ctx context.Context, userID uint) (*domain.Profile, error) {
	profile, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Dynamically evaluate verified & premium status via WordPress-style filter hook
	if u.hooks != nil {
		data := map[string]any{
			"user_id":     userID,
			"is_verified": profile.IsVerified,
			"is_premium":  profile.IsPremium,
		}
		filtered, err := u.hooks.ApplyFilter(ctx, "profile.decorate", data)
		if err == nil {
			if m, ok := filtered.(map[string]any); ok {
				if v, ok := m["is_verified"].(bool); ok {
					profile.IsVerified = v
				}
				if p, ok := m["is_premium"].(bool); ok {
					profile.IsPremium = p
				}
			}
		}
	}

	return profile, nil
}

func (u *profileUsecase) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	return u.repo.Create(ctx, profile)
}

func (u *profileUsecase) UpdateMyProfile(ctx context.Context, userID uint, req domain.UpdateProfileRequest) (*domain.Profile, error) {
	profile, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.Age != nil {
		profile.Age = *req.Age
	}
	if req.Bio != nil {
		profile.Bio = *req.Bio
	}
	if req.Gender != nil {
		profile.Gender = domain.Gender(*req.Gender)
	}
	if req.Location != nil {
		profile.Location = *req.Location
	}
	if req.Interests != nil {
		profile.Interests = req.Interests
	}
	if req.Photos != nil {
		profile.Photos = req.Photos
	}

	if err := u.repo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return profile, nil
}

func (u *profileUsecase) GetRandomProfiles(ctx context.Context, userID uint, filter domain.DiscoveryFilter) ([]domain.Profile, error) {
	// 1. Get daily viewed profiles count for current user
	viewCount, err := u.repo.GetDailyViewCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily view count: %w", err)
	}

	// 2. Check quota filter hook: "swipe.check_quota"
	hasUnlimited := false
	if u.hooks != nil {
		quotaData := map[string]any{
			"user_id":       userID,
			"has_unlimited": false,
		}
		res, err := u.hooks.ApplyFilter(ctx, "swipe.check_quota", quotaData)
		if err == nil {
			if m, ok := res.(map[string]any); ok {
				if val, ok := m["has_unlimited"].(bool); ok {
					hasUnlimited = val
				}
			}
		}
	}

	if !hasUnlimited && viewCount >= u.limitSwipe {
		return nil, errors.New("daily swipe limit reached. Upgrade to premium for unlimited swipes")
	}

	// Determine remaining limit for free users
	if !hasUnlimited {
		remaining := u.limitSwipe - viewCount
		if filter.Limit > remaining || filter.Limit <= 0 {
			filter.Limit = remaining
		}
	} else if filter.Limit <= 0 {
		filter.Limit = 10
	}

	profiles, err := u.repo.GetRandomProfiles(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch discovery profiles: %w", err)
	}

	return profiles, nil
}
