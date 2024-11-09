package services

import (
	"fmt"

	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"
)

type ProfileService interface {
	GetProfileByID(id string) (*models.Profile, error)
}

type profileService struct {
	profileRepo      repositories.ProfileRepository
	subscriptionRepo repositories.SubscriptionRepository
}

func NewProfileService(profileRepo repositories.ProfileRepository, subscriptionRepo repositories.SubscriptionRepository) *profileService {
	return &profileService{
		profileRepo:      profileRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *profileService) GetProfileByID(id string) (*models.Profile, error) {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not get profile: %w", err)
	}

	// Check if user has an active subscription
	profile.IsPremium, err = s.subscriptionRepo.GetActiveSubscription(profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not get active subscription: %w", err)
	}

	// Check if user has verified label subscription
	hasVerifiedLabel, err := s.subscriptionRepo.HasFeature(profile.UserID, models.FeatureNameVerifiedLabel)
	if err != nil {
		return nil, err
	}

	profile.VerifiedLabel = hasVerifiedLabel
	return profile, nil
}
