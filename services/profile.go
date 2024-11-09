package services

import (
	"fmt"

	model "kopelko-dating-app-backend/models"
	repository "kopelko-dating-app-backend/repositories"
)

type ProfileService interface {
	GetProfileByID(id string) (*model.Profile, error)
}

type profileService struct {
	profileRepo      repository.ProfileRepository
	subscriptionRepo repository.SubscriptionRepository
}

func NewProfileService(profileRepo repository.ProfileRepository, subscriptionRepo repository.SubscriptionRepository) *profileService {
	return &profileService{
		profileRepo:      profileRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *profileService) GetProfileByID(id string) (*model.Profile, error) {
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
	hasVerifiedLabel, err := s.subscriptionRepo.HasFeature(profile.UserID, model.FeatureNameVerifiedLabel)
	if err != nil {
		return nil, err
	}

	profile.VerifiedLabel = hasVerifiedLabel
	return profile, nil
}
