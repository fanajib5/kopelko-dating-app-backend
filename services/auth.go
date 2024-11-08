package services

import (
	"kopelko-dating-app-backend/dto"
	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	profileRepo *repositories.ProfileRepository
}

func NewAuthService(userRepo *repositories.UserRepository, profileRepo *repositories.ProfileRepository) *AuthService {
	return &AuthService{userRepo: userRepo, profileRepo: profileRepo}
}

func (s *AuthService) RegisterUser(req *dto.RegisterRequest) (*models.User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx := s.userRepo.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPwd),
	}

	if err := s.userRepo.CreateUserTx(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	profile := &models.Profile{
		UserID:   user.ID,
		Name:     req.Name,
		Age:      req.Age,
		Bio:      req.Bio,
		Gender:   req.Gender,
		Location: req.Location,
	}

	if err := s.profileRepo.CreateProfileTx(tx, profile); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, nil
}
