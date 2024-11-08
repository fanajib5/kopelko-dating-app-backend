package services

import (
	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"
)

type ProfileService interface {
	GetProfileByID(id string) (*models.Profile, error)
}

type profileService struct {
	profileRepo repositories.ProfileRepository
}

func NewProfileService(profileRepo repositories.ProfileRepository) *profileService {
	return &profileService{profileRepo: profileRepo}
}

func (s *profileService) GetProfileByID(id string) (*models.Profile, error) {
	return s.profileRepo.FindByID(id)
}
