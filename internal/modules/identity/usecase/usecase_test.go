package usecase_test

import (
	"context"
	"errors"
	"testing"

	"kopelko-dating-app-backend/internal/modules/identity/domain"
	"kopelko-dating-app-backend/internal/modules/identity/usecase"
	profileDomain "kopelko-dating-app-backend/internal/modules/profile/domain"
	"kopelko-dating-app-backend/internal/platform/database"
	"kopelko-dating-app-backend/internal/platform/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *MockUserRepository) CreateWithTx(ctx context.Context, tx database.DBTX, user *domain.User) error {
	args := m.Called(ctx, tx, user)
	if args.Error(0) == nil {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) SetVerified(ctx context.Context, id uint, isVerified bool) error {
	args := m.Called(ctx, id, isVerified)
	return args.Error(0)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateToken(tokenStr string) (*token.JWTCustomClaims, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*token.JWTCustomClaims), args.Error(1)
}

type MockProfileRepoForIdentity struct {
	mock.Mock
}

func (m *MockProfileRepoForIdentity) Create(ctx context.Context, profile *profileDomain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepoForIdentity) CreateWithTx(ctx context.Context, tx database.DBTX, profile *profileDomain.Profile) error {
	args := m.Called(ctx, tx, profile)
	return args.Error(0)
}

func (m *MockProfileRepoForIdentity) GetByUserID(ctx context.Context, userID uint) (*profileDomain.Profile, error) {
	return nil, nil
}
func (m *MockProfileRepoForIdentity) Update(ctx context.Context, profile *profileDomain.Profile) error {
	return nil
}
func (m *MockProfileRepoForIdentity) RecordView(ctx context.Context, userID, viewedUserID uint, swipeID *uint) error {
	return nil
}
func (m *MockProfileRepoForIdentity) RecordViewWithTx(ctx context.Context, tx database.DBTX, userID, viewedUserID uint, swipeID *uint) error {
	return nil
}
func (m *MockProfileRepoForIdentity) GetDailyViewCount(ctx context.Context, userID uint) (int, error) {
	return 0, nil
}
func (m *MockProfileRepoForIdentity) GetRandomProfiles(ctx context.Context, currentUserID uint, limit int) ([]profileDomain.Profile, error) {
	return nil, nil
}

func TestIdentityUsecase_Register_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepoForIdentity)
	tokenSvc := new(MockTokenService)
	svc := usecase.NewIdentityUsecase(userRepo, profileRepo, tokenSvc, nil)

	userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Profile")).Return(nil)
	tokenSvc.On("GenerateToken", uint(1), "john@example.com").Return("dummy-jwt-token", nil)

	user, tokenStr, err := svc.Register(
		context.Background(),
		"john@example.com", "password123", "John Doe",
		25, "male", "Surabaya", []string{"coding"}, []string{"https://img.com/1.jpg"},
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "dummy-jwt-token", tokenStr)
	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
}

func TestIdentityUsecase_Login_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepoForIdentity)
	tokenSvc := new(MockTokenService)
	svc := usecase.NewIdentityUsecase(userRepo, profileRepo, tokenSvc, nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:           1,
		Email:        "user@test.com",
		PasswordHash: string(hashedPassword),
	}

	userRepo.On("GetByEmail", mock.Anything, "user@test.com").Return(existingUser, nil)
	tokenSvc.On("GenerateToken", uint(1), "user@test.com").Return("jwt-token-123", nil)

	user, tokenStr, err := svc.Login(context.Background(), "user@test.com", "secret123")

	assert.NoError(t, err)
	assert.Equal(t, existingUser, user)
	assert.Equal(t, "jwt-token-123", tokenStr)
	userRepo.AssertExpectations(t)
}

func TestIdentityUsecase_Login_WrongPassword(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepoForIdentity)
	tokenSvc := new(MockTokenService)
	svc := usecase.NewIdentityUsecase(userRepo, profileRepo, tokenSvc, nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:           1,
		Email:        "user@test.com",
		PasswordHash: string(hashedPassword),
	}

	userRepo.On("GetByEmail", mock.Anything, "user@test.com").Return(existingUser, nil)

	user, tokenStr, err := svc.Login(context.Background(), "user@test.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, tokenStr)
	assert.Equal(t, "invalid email or password", err.Error())
}
