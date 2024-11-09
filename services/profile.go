package services

import (
	"fmt"
	"time"

	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"
)

type ProfileService interface {
	GetProfileByID(id uint) (*models.Profile, error)
	GetRandomProfiles(viewerID uint) (*models.Profile, error)
}

type profileService struct {
	profileRepo      repositories.ProfileRepository
	profileViewRepo  repositories.ProfileViewRepository
	subscriptionRepo repositories.SubscriptionRepository
	limitView        int
}

func NewProfileService(profileRepo repositories.ProfileRepository, profileViewRepo repositories.ProfileViewRepository, subscriptionRepo repositories.SubscriptionRepository, limitView int) *profileService {
	return &profileService{
		profileRepo:      profileRepo,
		profileViewRepo:  profileViewRepo,
		subscriptionRepo: subscriptionRepo,
		limitView:        limitView,
	}
}

func (s *profileService) GetProfileByID(id uint) (*models.Profile, error) {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not get profile: %w", err)
	}

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

func (s *profileService) GetRandomProfiles(viewerID uint) (*models.Profile, error) {
	profile, err := s.profileViewRepo.GetUnviewedProfiles(viewerID)
	if err != nil {
		return nil, fmt.Errorf("could not get unviewed profiles: %w", err)
	}

	// If there are no profiles to view, return an error
	if profile == nil {
		return nil, fmt.Errorf("no profiles to view, you've reached the limit to view profiles today")
	}

	// Log views in the `profile_views` table
	view := models.ProfileView{
		UserID:       viewerID,
		ViewedUserID: profile.UserID,
		ViewDate:     time.Now(),
	}
	if err := s.profileViewRepo.CreateProfileView(&view); err != nil {
		return nil, fmt.Errorf("could not create profile view: %w", err)
	}

	return profile, nil
}
