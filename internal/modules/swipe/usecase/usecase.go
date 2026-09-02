package usecase

import (
	"context"
	"errors"
	"fmt"

	profileDomain "kopelko-dating-app-backend/internal/modules/profile/domain"
	subscriptionDomain "kopelko-dating-app-backend/internal/modules/subscription/domain"
	"kopelko-dating-app-backend/internal/modules/swipe/domain"
	"kopelko-dating-app-backend/internal/platform/database"

	"github.com/jackc/pgx/v5"
)

type swipeUsecase struct {
	swipeRepo       domain.SwipeRepository
	subscriptionSvc subscriptionDomain.SubscriptionService
	profileRepo     profileDomain.ProfileRepository
	transactor      database.Transactor
	limitSwipe      int
}

func NewSwipeUsecase(
	swipeRepo domain.SwipeRepository,
	subSvc subscriptionDomain.SubscriptionService,
	profileRepo profileDomain.ProfileRepository,
	transactor database.Transactor,
	limitSwipe int,
) domain.SwipeService {
	return &swipeUsecase{
		swipeRepo:       swipeRepo,
		subscriptionSvc: subSvc,
		profileRepo:     profileRepo,
		transactor:      transactor,
		limitSwipe:      limitSwipe,
	}
}

func (u *swipeUsecase) SwipeUser(ctx context.Context, userID, targetUserID uint, swipeType domain.SwipeType) (*domain.SwipeResponse, error) {
	if userID == targetUserID {
		return nil, errors.New("cannot swipe yourself")
	}

	// Check if user has unlimited quota
	hasNoQuota, err := u.subscriptionSvc.HasActiveFeature(ctx, userID, "no_swipe_quota")
	if err != nil {
		return nil, fmt.Errorf("failed to verify quota feature: %w", err)
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
		if tx != nil {
			_ = u.profileRepo.RecordViewWithTx(ctx, tx, userID, targetUserID, &swipe.ID)
		} else {
			_ = u.profileRepo.RecordView(ctx, userID, targetUserID, &swipe.ID)
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

		return nil
	}

	if u.transactor != nil {
		err = u.transactor.WithinTransaction(ctx, func(tx pgx.Tx) error {
			return execSwipe(tx)
		})
	} else {
		err = execSwipe(nil)
	}

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (u *swipeUsecase) GetMatches(ctx context.Context, userID uint) ([]domain.MatchDetail, error) {
	return u.swipeRepo.GetMatchesByUserID(ctx, userID)
}
