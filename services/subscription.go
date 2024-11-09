package services

import (
	"errors"
	"fmt"
	"time"

	model "kopelko-dating-app-backend/models"
	repository "kopelko-dating-app-backend/repositories"
)

type SubscriptionService interface {
	SubscribeUser(userID uint, featureID int) error
}

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
	featureRepo      repository.PremiumFeatureRepository
}

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository, featureRepo repository.PremiumFeatureRepository) *subscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		featureRepo:      featureRepo,
	}
}

// SubscribeUser subscribes a user to a premium feature
func (s *subscriptionService) SubscribeUser(userID uint, featureID int) error {
	feature, err := s.featureRepo.GetFeatureByID(featureID)
	if feature == nil {
		return errors.New("premium feature not found")
	}
	if err != nil {
		return fmt.Errorf("could not get premium feature: %w", err)
	}

	if featureID < 0 {
		return errors.New("invalid feature ID")
	}
	featureIDuint := uint(featureID)

	now := time.Now()

	// Set subscription duration (e.g., 1 month)
	startDate := now
	endDate := now.AddDate(0, 1, 0)

	subscription := &model.Subscription{
		UserID:    userID,
		FeatureID: featureIDuint,
		StartDate: startDate,
		EndDate:   endDate,
		AutoRenew: true,
	}

	err = s.subscriptionRepo.CreateSubscription(subscription)
	if err != nil {
		return fmt.Errorf("could not create subscription: %w", err)
	}

	return nil
}
