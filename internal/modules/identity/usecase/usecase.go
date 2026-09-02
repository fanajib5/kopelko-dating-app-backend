package usecase

import (
	"context"
	"errors"
	"fmt"

	"kopelko-dating-app-backend/internal/modules/identity/domain"
	profileDomain "kopelko-dating-app-backend/internal/modules/profile/domain"
	"kopelko-dating-app-backend/internal/platform/database"
	"kopelko-dating-app-backend/internal/platform/token"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type identityUsecase struct {
	userRepo    domain.UserRepository
	profileRepo profileDomain.ProfileRepository
	tokenSvc    token.TokenService
	transactor  database.Transactor
}

func NewIdentityUsecase(
	userRepo domain.UserRepository,
	profileRepo profileDomain.ProfileRepository,
	tokenSvc token.TokenService,
	transactor database.Transactor,
) domain.IdentityService {
	return &identityUsecase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		tokenSvc:    tokenSvc,
		transactor:  transactor,
	}
}

func (u *identityUsecase) Register(
	ctx context.Context,
	email, password, name string,
	age int,
	gender, location string,
	interests, photos []string,
) (*domain.User, string, error) {
	// Check existing user
	existing, _ := u.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", errors.New("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsVerified:   false,
	}

	profile := &profileDomain.Profile{
		Name:      name,
		Age:       age,
		Bio:       "",
		Gender:    profileDomain.Gender(gender),
		Location:  location,
		Interests: interests,
		Photos:    photos,
		IsPremium: false,
	}

	// Atomic registration transaction: create user + create profile
	if u.transactor != nil {
		err = u.transactor.WithinTransaction(ctx, func(tx pgx.Tx) error {
			if err := u.userRepo.CreateWithTx(ctx, tx, user); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			profile.UserID = user.ID
			if err := u.profileRepo.CreateWithTx(ctx, tx, profile); err != nil {
				return fmt.Errorf("failed to create profile: %w", err)
			}
			return nil
		})
	} else {
		// Fallback without transactor (e.g. in tests with plain mocks)
		if err := u.userRepo.Create(ctx, user); err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}
		profile.UserID = user.ID
		if err := u.profileRepo.Create(ctx, profile); err != nil {
			return nil, "", fmt.Errorf("failed to create profile: %w", err)
		}
	}

	if err != nil {
		return nil, "", err
	}

	tokenStr, err := u.tokenSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, tokenStr, nil
}

func (u *identityUsecase) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	tokenStr, err := u.tokenSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, tokenStr, nil
}
