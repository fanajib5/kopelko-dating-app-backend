package services

import (
	"fmt"
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
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	tx := s.userRepo.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("could not start transaction: %w", tx.Error)
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPwd),
	}

	if err := s.userRepo.CreateUserTx(tx, user); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("could not create user: %w", err)
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
		return nil, fmt.Errorf("could not create profile: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return user, nil
}
