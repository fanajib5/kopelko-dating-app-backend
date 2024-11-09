package services

import (
	"fmt"

	"kopelko-dating-app-backend/dto"
	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/repositories"
	"kopelko-dating-app-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(req *dto.RegisterRequest) (*models.User, error)
	LoginUser(req *dto.LoginRequest) (*models.User, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	profileRepo repositories.ProfileRepository
}

func NewAuthService(userRepo repositories.UserRepository, profileRepo repositories.ProfileRepository) *authService {
	return &authService{userRepo: userRepo, profileRepo: profileRepo}
}

func (s *authService) RegisterUser(req *dto.RegisterRequest) (*models.User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	tx := s.userRepo.BeginTx()
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

func (s *authService) LoginUser(req *dto.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Assuming you have a method to generate a token for the user
	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	user.Token = token
	return user, nil
}

func (s *authService) generateToken(user *models.User) (string, error) {
	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}
	return token, nil
}
