package services

import (
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
	return s.profileRepo.FindByID(id)
}
