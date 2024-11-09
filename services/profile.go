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
	profileRepo repository.ProfileRepository
}

func NewProfileService(profileRepo repository.ProfileRepository) *profileService {
	return &profileService{profileRepo: profileRepo}
}

func (s *profileService) GetProfileByID(id string) (*model.Profile, error) {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not get profile: %w", err)
	}
	return profile, nil
}
