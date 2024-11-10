package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"

	"github.com/labstack/echo/v4"
)

type SwipeService interface {
	SwipeProfile(userID uint, targetUserID int, swipeType string) error
}

type swipeService struct {
	swipeRepo        repositories.SwipeRepository
	subscriptionRepo repositories.SubscriptionRepository
	profileViewRepo  repositories.ProfileViewRepository
	maxSwipes        int
}

// NewSwipeService creates a new SwipeService with a maximum swipe limit
func NewSwipeService(swipeRepo repositories.SwipeRepository, subscriptionRepo repositories.SubscriptionRepository, profileViewRepo repositories.ProfileViewRepository, maxSwipes int) *swipeService {
	return &swipeService{
		swipeRepo:        swipeRepo,
		subscriptionRepo: subscriptionRepo,
		profileViewRepo:  profileViewRepo,
		maxSwipes:        maxSwipes,
	}
}

// SwipeProfile handles the swipe logic with daily limits and swipe uniqueness
func (s *swipeService) SwipeProfile(userID uint, targetUserID int, swipeType string) error {
	now := time.Now()

	log.Printf("Attempting to swipe an user with userID: %d, targetUserID: %d, swipeType: %s", userID, targetUserID, swipeType)

	if targetUserID < 0 {
		return errors.New("invalid target user ID")
	}
	targetUserIDuint := uint(targetUserID)

	// Check if the user has an active subscription with "no_swipe_quota" feature
	hasUnlimitedSwipes, err := s.subscriptionRepo.HasFeature(userID, models.FeatureNameNoSwipeQuota)
	if err != nil {
		return fmt.Errorf("could not check user subscription: %w", err)
	}

	// Allow unlimited swipes if user is premium
	if !hasUnlimitedSwipes {
		if err := s.checkDailySwipes(userID, now); err != nil {
			return fmt.Errorf("could not check daily swipes: %w", err)
		}
	}

	// Check if the user has already swiped on this target today
	hasSwiped, err := s.swipeRepo.HasSwipedToday(userID, targetUserIDuint, now.Truncate(24*time.Hour))
	if err != nil {
		return fmt.Errorf("could not check if user has swiped: %w", err)
	}
	if hasSwiped {
		return echo.NewHTTPError(http.StatusConflict, "Already swiped on this user today")
	}

	log.Println("Creating swipe and profile view data")
	// Set up the Swipe data
	swipe := models.Swipe{
		UserID:       userID,
		TargetUserID: targetUserIDuint,
		SwipeType:    swipeType,
		SwipeDate:    now,
	}

	// Set up the ProfileView data
	profileView := models.ProfileView{
		UserID:       userID,
		ViewedUserID: targetUserIDuint,
		ViewDate:     now,
	}

	swipeAndVeiw := models.SwipeAndViewData{
		Swipe:       swipe,
		ProfileView: profileView,
	}

	if err = s.profileViewRepo.CreateSwipeAndView(swipeAndVeiw); err != nil {
		return fmt.Errorf("could not create swipe: %w", err)
	}

	return nil
}

// Check daily swipe count
func (s *swipeService) checkDailySwipes(userID uint, now time.Time) error {
	log.Println("Checking daily swipes limit for user:", userID)
	today := now.Truncate(24 * time.Hour)
	count, err := s.swipeRepo.GetDailySwipes(userID, today)
	if err != nil {
		return fmt.Errorf("could not get daily swipes: %w", err)
	}
	if count >= int64(s.maxSwipes) {
		return echo.NewHTTPError(http.StatusForbidden, "Swipe limit reached for today")
	}
	return nil
}
