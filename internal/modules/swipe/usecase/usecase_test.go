package usecase_test

import (
	"context"
	"testing"
	"time"

	profileDomain "kopelko-dating-app-backend/internal/modules/profile/domain"
	subscriptionDomain "kopelko-dating-app-backend/internal/modules/subscription/domain"
	"kopelko-dating-app-backend/internal/modules/swipe/domain"
	"kopelko-dating-app-backend/internal/modules/swipe/usecase"
	"kopelko-dating-app-backend/internal/platform/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSwipeRepository struct {
	mock.Mock
}

func (m *MockSwipeRepository) CreateSwipe(ctx context.Context, swipe *domain.Swipe) error {
	args := m.Called(ctx, swipe)
	return args.Error(0)
}

func (m *MockSwipeRepository) CreateSwipeWithTx(ctx context.Context, tx database.DBTX, swipe *domain.Swipe) error {
	args := m.Called(ctx, tx, swipe)
	return args.Error(0)
}

func (m *MockSwipeRepository) GetDailySwipeCount(ctx context.Context, userID uint) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockSwipeRepository) GetDailySwipeCountWithTx(ctx context.Context, tx database.DBTX, userID uint) (int, error) {
	args := m.Called(ctx, tx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockSwipeRepository) HasSwiped(ctx context.Context, userID, targetUserID uint) (bool, error) {
	args := m.Called(ctx, userID, targetUserID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSwipeRepository) HasSwipedWithTx(ctx context.Context, tx database.DBTX, userID, targetUserID uint) (bool, error) {
	args := m.Called(ctx, tx, userID, targetUserID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSwipeRepository) GetSwipe(ctx context.Context, userID, targetUserID uint) (*domain.Swipe, error) {
	args := m.Called(ctx, userID, targetUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Swipe), args.Error(1)
}

func (m *MockSwipeRepository) GetSwipeWithTx(ctx context.Context, tx database.DBTX, userID, targetUserID uint) (*domain.Swipe, error) {
	args := m.Called(ctx, tx, userID, targetUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Swipe), args.Error(1)
}

func (m *MockSwipeRepository) CreateMatch(ctx context.Context, user1ID, user2ID uint) (*domain.Match, error) {
	args := m.Called(ctx, user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Match), args.Error(1)
}

func (m *MockSwipeRepository) CreateMatchWithTx(ctx context.Context, tx database.DBTX, user1ID, user2ID uint) (*domain.Match, error) {
	args := m.Called(ctx, tx, user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Match), args.Error(1)
}

func (m *MockSwipeRepository) GetMatch(ctx context.Context, user1ID, user2ID uint) (*domain.Match, error) {
	args := m.Called(ctx, user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Match), args.Error(1)
}

func (m *MockSwipeRepository) GetMatchesByUserID(ctx context.Context, userID uint) ([]domain.MatchDetail, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MatchDetail), args.Error(1)
}

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) HasActiveFeature(ctx context.Context, userID uint, featureName string) (bool, error) {
	args := m.Called(ctx, userID, featureName)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionService) Subscribe(ctx context.Context, userID uint, featureName string) (*subscriptionDomain.Subscription, error) {
	return nil, nil
}

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) Create(ctx context.Context, profile *profileDomain.Profile) error {
	return nil
}
func (m *MockProfileRepository) CreateWithTx(ctx context.Context, tx database.DBTX, profile *profileDomain.Profile) error {
	return nil
}
func (m *MockProfileRepository) GetByUserID(ctx context.Context, userID uint) (*profileDomain.Profile, error) {
	return nil, nil
}
func (m *MockProfileRepository) Update(ctx context.Context, profile *profileDomain.Profile) error {
	return nil
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
	return 0, nil
}
func (m *MockProfileRepository) GetRandomProfiles(ctx context.Context, currentUserID uint, filter profileDomain.DiscoveryFilter) ([]profileDomain.Profile, error) {
	return nil, nil
}

func TestSwipeUsecase_SwipeLike_MutualMatch(t *testing.T) {
	swipeRepo := new(MockSwipeRepository)
	subSvc := new(MockSubscriptionService)
	profileRepo := new(MockProfileRepository)
	svc := usecase.NewSwipeUsecase(swipeRepo, subSvc, profileRepo, nil, 10)

	subSvc.On("HasActiveFeature", mock.Anything, uint(1), "no_swipe_quota").Return(false, nil)
	swipeRepo.On("HasSwiped", mock.Anything, uint(1), uint(2)).Return(false, nil)
	swipeRepo.On("GetDailySwipeCount", mock.Anything, uint(1)).Return(3, nil)
	swipeRepo.On("CreateSwipe", mock.Anything, mock.AnythingOfType("*domain.Swipe")).Return(nil)
	profileRepo.On("RecordView", mock.Anything, uint(1), uint(2), mock.Anything).Return(nil)

	targetSwipe := &domain.Swipe{UserID: 2, TargetUserID: 1, SwipeType: domain.SwipeLike}
	swipeRepo.On("GetSwipe", mock.Anything, uint(2), uint(1)).Return(targetSwipe, nil)
	match := &domain.Match{ID: 100, User1ID: 1, User2ID: 2}
	swipeRepo.On("CreateMatch", mock.Anything, uint(1), uint(2)).Return(match, nil)

	res, err := svc.SwipeUser(context.Background(), 1, 2, domain.SwipeLike)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.IsMatch)
	assert.Equal(t, match, res.Match)
	swipeRepo.AssertExpectations(t)
}

func TestSwipeUsecase_GetMatches(t *testing.T) {
	swipeRepo := new(MockSwipeRepository)
	subSvc := new(MockSubscriptionService)
	profileRepo := new(MockProfileRepository)
	svc := usecase.NewSwipeUsecase(swipeRepo, subSvc, profileRepo, nil, 10)

	expectedMatches := []domain.MatchDetail{
		{
			MatchID:       1,
			MatchedUserID: 2,
			MatchedAt:     time.Now(),
			Profile:       &profileDomain.Profile{ID: 2, UserID: 2, Name: "Jane"},
		},
	}
	swipeRepo.On("GetMatchesByUserID", mock.Anything, uint(1)).Return(expectedMatches, nil)

	matches, err := svc.GetMatches(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, matches, 1)
	assert.Equal(t, "Jane", matches[0].Profile.Name)
	swipeRepo.AssertExpectations(t)
}
