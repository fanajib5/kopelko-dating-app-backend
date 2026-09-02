package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"kopelko-dating-app-backend/internal/plugins/subscription/domain"
	"kopelko-dating-app-backend/internal/plugins/subscription/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) GetFeatureByName(ctx context.Context, featureName string) (*domain.PremiumFeature, error) {
	args := m.Called(ctx, featureName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PremiumFeature), args.Error(1)
}

func (m *MockSubscriptionRepository) GetFeatureByID(ctx context.Context, featureID uint) (*domain.PremiumFeature, error) {
	args := m.Called(ctx, featureID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PremiumFeature), args.Error(1)
}

func (m *MockSubscriptionRepository) CreateOrRenewSubscription(ctx context.Context, sub *domain.Subscription) error {
	args := m.Called(ctx, sub)
	if args.Error(0) == nil {
		sub.ID = 1
	}
	return args.Error(0)
}

func (m *MockSubscriptionRepository) HasActiveFeature(ctx context.Context, userID uint, featureName string) (bool, error) {
	args := m.Called(ctx, userID, featureName)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionRepository) GetActiveSubscriptions(ctx context.Context, userID uint) ([]domain.Subscription, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Subscription), args.Error(1)
}

func TestSubscriptionUsecase_Subscribe_Success(t *testing.T) {
	repo := new(MockSubscriptionRepository)
	svc := usecase.NewSubscriptionUsecase(repo)

	feature := &domain.PremiumFeature{ID: 1, FeatureName: "no_swipe_quota", Description: "Unlimited swipes"}
	repo.On("GetFeatureByName", mock.Anything, "no_swipe_quota").Return(feature, nil)
	repo.On("CreateOrRenewSubscription", mock.Anything, mock.AnythingOfType("*domain.Subscription")).Return(nil)

	sub, err := svc.Subscribe(context.Background(), 10, "no_swipe_quota")

	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, uint(10), sub.UserID)
	assert.Equal(t, uint(1), sub.FeatureID)
	assert.True(t, sub.EndDate.After(time.Now()))
	repo.AssertExpectations(t)
}

func TestSubscriptionUsecase_Subscribe_FeatureNotFound(t *testing.T) {
	repo := new(MockSubscriptionRepository)
	svc := usecase.NewSubscriptionUsecase(repo)

	repo.On("GetFeatureByName", mock.Anything, "invalid_feature").Return(nil, errors.New("not found"))

	sub, err := svc.Subscribe(context.Background(), 10, "invalid_feature")

	assert.Error(t, err)
	assert.Nil(t, sub)
	repo.AssertExpectations(t)
}

func TestSubscriptionUsecase_HasActiveFeature(t *testing.T) {
	repo := new(MockSubscriptionRepository)
	svc := usecase.NewSubscriptionUsecase(repo)

	repo.On("HasActiveFeature", mock.Anything, uint(10), "verified_label").Return(true, nil)

	hasActive, err := svc.HasActiveFeature(context.Background(), 10, "verified_label")

	assert.NoError(t, err)
	assert.True(t, hasActive)
	repo.AssertExpectations(t)
}
