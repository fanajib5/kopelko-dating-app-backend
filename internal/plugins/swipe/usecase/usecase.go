package usecase

import (
	"context"
	"errors"
	"fmt"

	"kopelko-dating-app-backend/internal/core/hook"
	"kopelko-dating-app-backend/internal/platform/database"
	profileDomain "kopelko-dating-app-backend/internal/plugins/profile/domain"
	"kopelko-dating-app-backend/internal/plugins/swipe/domain"

	"github.com/jackc/pgx/v5"
)

type swipeUsecase struct {
	swipeRepo   domain.SwipeRepository
	profileRepo profileDomain.ProfileRepository
	transactor  database.Transactor
	hooks       hook.HookManager
	limitSwipe  int
}

func NewSwipeUsecase(
	swipeRepo domain.SwipeRepository,
	profileRepo profileDomain.ProfileRepository,
	transactor database.Transactor,
	hooks hook.HookManager,
	limitSwipe int,
) domain.SwipeService {
	return &swipeUsecase{
		swipeRepo:   swipeRepo,
		profileRepo: profileRepo,
		transactor:  transactor,
		hooks:       hooks,
		limitSwipe:  limitSwipe,
	}
}

func (u *swipeUsecase) SwipeUser(ctx context.Context, userID, targetUserID uint, swipeType domain.SwipeType) (*domain.SwipeResponse, error) {
	if userID == targetUserID {
		return nil, errors.New("cannot swipe yourself")
	}

	// Check if user has unlimited quota via Filter Hook: "swipe.check_quota"
	hasNoQuota := false
	if u.hooks != nil {
		quotaData := map[string]any{
			"user_id":       userID,
			"has_unlimited": false,
		}
		res, err := u.hooks.ApplyFilter(ctx, "swipe.check_quota", quotaData)
		if err == nil {
			if m, ok := res.(map[string]any); ok {
				if val, ok := m["has_unlimited"].(bool); ok {
					hasNoQuota = val
				}
			}
		}
	}

	var resp domain.SwipeResponse

	execSwipe := func(tx database.DBTX) error {
		// 1. Check if already swiped
		var alreadySwiped bool
		var err error
		if tx != nil {
			alreadySwiped, err = u.swipeRepo.HasSwipedWithTx(ctx, tx, userID, targetUserID)
		} else {
			alreadySwiped, err = u.swipeRepo.HasSwiped(ctx, userID, targetUserID)
		}
		if err != nil {
			return fmt.Errorf("failed to check existing swipe: %w", err)
		}
		if alreadySwiped {
			return errors.New("you have already swiped on this profile")
		}

		// 2. Check quota inside transaction
		if !hasNoQuota {
			var count int
			if tx != nil {
				count, err = u.swipeRepo.GetDailySwipeCountWithTx(ctx, tx, userID)
			} else {
				count, err = u.swipeRepo.GetDailySwipeCount(ctx, userID)
			}
			if err != nil {
				return fmt.Errorf("failed to count daily swipes: %w", err)
			}
			if count >= u.limitSwipe {
				return errors.New("daily swipe limit reached. Upgrade to premium for unlimited swipes")
			}
		}

		// 3. Create swipe record
		swipe := &domain.Swipe{
			UserID:       userID,
			TargetUserID: targetUserID,
			SwipeType:    swipeType,
		}
		if tx != nil {
			err = u.swipeRepo.CreateSwipeWithTx(ctx, tx, swipe)
		} else {
			err = u.swipeRepo.CreateSwipe(ctx, swipe)
		}
		if err != nil {
			return fmt.Errorf("failed to save swipe: %w", err)
		}

		// 4. Record profile view
		if u.profileRepo != nil {
			if tx != nil {
				_ = u.profileRepo.RecordViewWithTx(ctx, tx, userID, targetUserID, &swipe.ID)
			} else {
				_ = u.profileRepo.RecordView(ctx, userID, targetUserID, &swipe.ID)
			}
		}

		resp.Swipe = swipe
		resp.IsMatch = false

		// 5. Check mutual like (Match)
		if swipeType == domain.SwipeLike {
			var targetSwipe *domain.Swipe
			if tx != nil {
				targetSwipe, err = u.swipeRepo.GetSwipeWithTx(ctx, tx, targetUserID, userID)
			} else {
				targetSwipe, err = u.swipeRepo.GetSwipe(ctx, targetUserID, userID)
			}

			if err == nil && targetSwipe != nil && targetSwipe.SwipeType == domain.SwipeLike {
				var match *domain.Match
				if tx != nil {
					match, err = u.swipeRepo.CreateMatchWithTx(ctx, tx, userID, targetUserID)
				} else {
					match, err = u.swipeRepo.CreateMatch(ctx, userID, targetUserID)
				}
				if err == nil {
					resp.IsMatch = true
					resp.Match = match
				}
			}
		}

		// Trigger Action Hook: "swipe.created"
		if u.hooks != nil {
			_ = u.hooks.DoAction(ctx, "swipe.created", &resp)
		}

		return nil
	}

	if u.transactor != nil {
		err := u.transactor.WithinTransaction(ctx, func(tx pgx.Tx) error {
			return execSwipe(tx)
		})
		if err != nil {
			return nil, err
		}
	} else {
		if err := execSwipe(nil); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}

func (u *swipeUsecase) GetMatches(ctx context.Context, userID uint) ([]domain.MatchDetail, error) {
	return u.swipeRepo.GetMatchesByUserID(ctx, userID)
}
