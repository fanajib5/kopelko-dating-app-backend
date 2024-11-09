package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	model "kopelko-dating-app-backend/models"
	repository "kopelko-dating-app-backend/repositories"

	"github.com/labstack/echo/v4"
)

type SwipeService interface {
	ProcessSwipe(userID uint, targetUserID int, swipeType string) error
}

type swipeService struct {
	swipeRepo        repository.SwipeRepository
	subscriptionRepo repository.SubscriptionRepository
	maxSwipes        int
}

// NewSwipeService creates a new SwipeService with a maximum swipe limit
func NewSwipeService(swipeRepo repository.SwipeRepository, subscriptionRepo repository.SubscriptionRepository, maxSwipes int) *swipeService {
	return &swipeService{
		swipeRepo:        swipeRepo,
		subscriptionRepo: subscriptionRepo,
		maxSwipes:        maxSwipes,
	}
}

// ProcessSwipe handles the swipe logic with daily limits and swipe uniqueness
func (s *swipeService) ProcessSwipe(userID uint, targetUserID int, swipeType string) error {
	now := time.Now()

	if targetUserID < 0 {
		return errors.New("invalid target user ID")
	}
	targetUserIDuint := uint(targetUserID)

	// Check if the user has an active subscription with "no_swipe_quota" feature
	hasUnlimitedSwipes, err := s.subscriptionRepo.HasFeature(userID, model.FeatureNameNoSwipeQuota)
	if err != nil {
		return err
	}

	// Allow unlimited swipes if user is premium
	if !hasUnlimitedSwipes {
		// Check daily swipe count
		today := now.Truncate(24 * time.Hour)
		count, err := s.swipeRepo.GetDailySwipes(userID, today)
		if err != nil {
			return fmt.Errorf("could not get daily swipes: %w", err)
		}
		if count >= int64(s.maxSwipes) {
			return echo.NewHTTPError(http.StatusForbidden, "Swipe limit reached for today")
		}
	}

	// Check if the user has already swiped on this target today
	hasSwiped, err := s.swipeRepo.HasSwipedToday(userID, targetUserIDuint, time.Now().Truncate(24*time.Hour))
	if err != nil {
		return fmt.Errorf("could not check if user has swiped: %w", err)
	}
	if hasSwiped {
		return echo.NewHTTPError(http.StatusConflict, "Already swiped on this user today")
	}

	// Create the swipe
	swipe := &model.Swipe{
		UserID:       userID,
		TargetUserID: targetUserIDuint,
		SwipeType:    swipeType,
		SwipeDate:    now.Truncate(24 * time.Hour),
	}
	return s.swipeRepo.CreateSwipe(swipe)
}
