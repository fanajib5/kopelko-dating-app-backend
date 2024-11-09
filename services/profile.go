package services

import (
	"fmt"

	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"
)

type ProfileService interface {
	GetProfileByID(id uint) (*models.Profile, error)
	GetRandomProfile() (*models.Profile, error)
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

func (s *profileService) GetProfileByID(id uint) (*models.Profile, error) {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not get profile: %w", err)
	}

	return s.getProfile(profile)
}

func (s *profileService) GetRandomProfile() (*models.Profile, error) {
	profile, err := s.profileRepo.FindRandom()
	if err != nil {
		return nil, fmt.Errorf("could not get random profile: %w", err)
	}

	return s.getProfile(profile)
}

func (s *profileService) getProfile(profile *models.Profile) (*models.Profile, error) {
	// Check if user has an active subscription
	subscription, err := s.subscriptionRepo.GetActiveSubscription(profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not get active subscription: %w", err)
	}
	if subscription != nil {
		profile.IsPremium = true
	}

	// Check if user has verified label subscription
	hasVerifiedLabel, err := s.subscriptionRepo.HasFeature(profile.UserID, models.FeatureNameVerifiedLabel)
	if err != nil {
		return nil, err
	}

	profile.VerifiedLabel = hasVerifiedLabel
	return profile, nil
}
