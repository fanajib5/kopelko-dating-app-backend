package usecase_test

import (
	"context"
	"testing"

	"kopelko-dating-app-backend/internal/modules/profile/domain"
	"kopelko-dating-app-backend/internal/modules/profile/usecase"
	subscriptionDomain "kopelko-dating-app-backend/internal/modules/subscription/domain"
	"kopelko-dating-app-backend/internal/platform/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepository) CreateWithTx(ctx context.Context, tx database.DBTX, profile *domain.Profile) error {
	args := m.Called(ctx, tx, profile)
	return args.Error(0)
}

func (m *MockProfileRepository) GetByUserID(ctx context.Context, userID uint) (*domain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepository) RecordView(ctx context.Context, userID, viewedUserID uint, swipeID *uint) error {
	args := m.Called(ctx, userID, viewedUserID, swipeID)
	return args.Error(0)
}

func (m *MockProfileRepository) RecordViewWithTx(ctx context.Context, tx database.DBTX, userID, viewedUserID uint, swipeID *uint) error {
	args := m.Called(ctx, tx, userID, viewedUserID, swipeID)
	return args.Error(0)
}

func (m *MockProfileRepository) GetDailyViewCount(ctx context.Context, userID uint) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockProfileRepository) GetRandomProfiles(ctx context.Context, currentUserID uint, limit int) ([]domain.Profile, error) {
	args := m.Called(ctx, currentUserID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Profile), args.Error(1)
}

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) Subscribe(ctx context.Context, userID uint, featureName string) (*subscriptionDomain.Subscription, error) {
	return nil, nil
}

func (m *MockSubscriptionService) HasActiveFeature(ctx context.Context, userID uint, featureName string) (bool, error) {
	args := m.Called(ctx, userID, featureName)
	return args.Bool(0), args.Error(1)
}

func TestProfileUsecase_GetMyProfile(t *testing.T) {
	repo := new(MockProfileRepository)
	subSvc := new(MockSubscriptionService)
	svc := usecase.NewProfileUsecase(repo, subSvc, 10)

	expected := &domain.Profile{ID: 1, UserID: 10, Name: "John"}
	repo.On("GetByUserID", mock.Anything, uint(10)).Return(expected, nil)

	p, err := svc.GetMyProfile(context.Background(), 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, p)
	repo.AssertExpectations(t)
}

func TestProfileUsecase_GetRandomProfiles_UnderQuota(t *testing.T) {
	repo := new(MockProfileRepository)
	subSvc := new(MockSubscriptionService)
	svc := usecase.NewProfileUsecase(repo, subSvc, 10)

	subSvc.On("HasActiveFeature", mock.Anything, uint(10), "no_swipe_quota").Return(false, nil)
	repo.On("GetDailyViewCount", mock.Anything, uint(10)).Return(2, nil)

	profiles := []domain.Profile{
		{ID: 1, UserID: 11, Name: "Alice"},
		{ID: 2, UserID: 12, Name: "Bob"},
	}
	repo.On("GetRandomProfiles", mock.Anything, uint(10), 8).Return(profiles, nil)
	repo.On("RecordView", mock.Anything, uint(10), uint(11), (*uint)(nil)).Return(nil)
	repo.On("RecordView", mock.Anything, uint(10), uint(12), (*uint)(nil)).Return(nil)

	res, err := svc.GetRandomProfiles(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	repo.AssertExpectations(t)
	subSvc.AssertExpectations(t)
}

func TestProfileUsecase_GetRandomProfiles_QuotaExceeded(t *testing.T) {
	repo := new(MockProfileRepository)
	subSvc := new(MockSubscriptionService)
	svc := usecase.NewProfileUsecase(repo, subSvc, 10)

	subSvc.On("HasActiveFeature", mock.Anything, uint(10), "no_swipe_quota").Return(false, nil)
	repo.On("GetDailyViewCount", mock.Anything, uint(10)).Return(10, nil)

	res, err := svc.GetRandomProfiles(context.Background(), 10)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "limit reached")
}

func TestProfileUsecase_GetRandomProfiles_UnlimitedPremium(t *testing.T) {
	repo := new(MockProfileRepository)
	subSvc := new(MockSubscriptionService)
	svc := usecase.NewProfileUsecase(repo, subSvc, 10)

	subSvc.On("HasActiveFeature", mock.Anything, uint(10), "no_swipe_quota").Return(true, nil)

	profiles := []domain.Profile{
		{ID: 1, UserID: 11, Name: "Alice"},
	}
	repo.On("GetRandomProfiles", mock.Anything, uint(10), 10).Return(profiles, nil)
	repo.On("RecordView", mock.Anything, uint(10), uint(11), (*uint)(nil)).Return(nil)

	res, err := svc.GetRandomProfiles(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	repo.AssertExpectations(t)
}
