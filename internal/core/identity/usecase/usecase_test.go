package usecase_test

import (
	"context"
	"errors"
	"testing"

	"kopelko-dating-app-backend/internal/core/hook"
	"kopelko-dating-app-backend/internal/core/identity/domain"
	"kopelko-dating-app-backend/internal/core/identity/usecase"
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

func TestIdentityUsecase_Register_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	tokenSvc := new(MockTokenService)
	hookMgr := hook.NewHookManager()

	var hookCalled bool
	hookMgr.AddAction("user.registered", 10, func(ctx context.Context, payload any) error {
		p, ok := payload.(*domain.UserRegisteredPayload)
		if ok && p.User.ID == 1 && p.Name == "John Doe" {
			hookCalled = true
		}
		return nil
	})

	svc := usecase.NewIdentityUsecase(userRepo, tokenSvc, nil, hookMgr)

	userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, errors.New("not found"))
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	tokenSvc.On("GenerateToken", uint(1), "john@example.com").Return("dummy-jwt-token", nil)

	user, tokenStr, err := svc.Register(
		context.Background(),
		"john@example.com", "password123", "John Doe",
		25, "male", "Surabaya", []string{"coding"}, []string{"https://img.com/1.jpg"},
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "dummy-jwt-token", tokenStr)
	assert.True(t, hookCalled, "user.registered hook should have been called")
	userRepo.AssertExpectations(t)
	tokenSvc.AssertExpectations(t)
}

func TestIdentityUsecase_Login_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	tokenSvc := new(MockTokenService)
	hookMgr := hook.NewHookManager()
	svc := usecase.NewIdentityUsecase(userRepo, tokenSvc, nil, hookMgr)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	mockUser := &domain.User{
		ID:           1,
		Email:        "john@example.com",
		PasswordHash: string(hashedPassword),
	}

	userRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(mockUser, nil)
	tokenSvc.On("GenerateToken", uint(1), "john@example.com").Return("dummy-jwt-token", nil)

	user, tokenStr, err := svc.Login(context.Background(), "john@example.com", "secret123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "dummy-jwt-token", tokenStr)
	userRepo.AssertExpectations(t)
}
