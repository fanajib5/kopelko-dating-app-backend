package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"
)

type SubscriptionService interface {
	SubscribeUser(userID uint, featureID int) error
}

type subscriptionService struct {
	subscriptionRepo repositories.SubscriptionRepository
	featureRepo      repositories.PremiumFeatureRepository
	profileRepo      repositories.ProfileRepository
}

func NewSubscriptionService(subscriptionRepo repositories.SubscriptionRepository, featureRepo repositories.PremiumFeatureRepository, profileRepo repositories.ProfileRepository) *subscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		featureRepo:      featureRepo,
		profileRepo:      profileRepo,
	}
}

// SubscribeUser subscribes a user to a premium feature
func (s *subscriptionService) SubscribeUser(userID uint, featureID int) error {
	log.Println("Subscire with userID: ", userID, " and featureID: ", featureID)

	feature, err := s.featureRepo.GetFeatureByID(featureID)
	if err != nil {
		return fmt.Errorf("could not get premium feature: %w", err)
	}
	if feature == nil {
		return errors.New("premium feature not found")
	}

	// Check if user has an active subscription
	log.Println("Checking if user has an active subscription")
	existingSubscription, err := s.subscriptionRepo.GetActiveSubscription(userID)
	if err != nil {
		return fmt.Errorf("could not get active subscription: %w", err)
	}
	if existingSubscription != nil {
		return errors.New("user already has an active subscription")
	}

	if featureID < 0 {
		return errors.New("invalid feature ID")
	}
	featureIDuint := uint(featureID)

	now := time.Now()

	log.Println("Creating subscription")
	// Set subscription duration (e.g., 1 month)
	// Assume that the subscription duration is 1 month
	startDate := now
	endDate := startDate.AddDate(0, 1, 0)

	subscription := &models.Subscription{
		UserID:    userID,
		FeatureID: featureIDuint,
		StartDate: startDate,
		EndDate:   endDate,
		AutoRenew: true,
	}

	tx := s.profileRepo.BeginTx()
	if tx.Error != nil {
		return fmt.Errorf("could not start transaction: %w", tx.Error)
	}

	if err = s.subscriptionRepo.CreateSubscriptionTx(tx, subscription); err != nil {
		tx.Rollback()
		return fmt.Errorf("could not create subscription: %w", err)
	}

	// Update profile's IsPremium to true
	if err := s.profileRepo.UpdateIsPremiumTx(tx, userID, true); err != nil {
		tx.Rollback()
		return fmt.Errorf("could not update profile: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
