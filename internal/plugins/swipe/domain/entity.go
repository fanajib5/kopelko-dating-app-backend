package domain

import (
	"context"
	"time"

	profileDomain "kopelko-dating-app-backend/internal/plugins/profile/domain"
	"kopelko-dating-app-backend/internal/platform/database"
)

type SwipeType string

const (
	SwipePass SwipeType = "pass"
	SwipeLike SwipeType = "like"
)

type Swipe struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	TargetUserID uint       `json:"target_user_id"`
	SwipeType    SwipeType  `json:"swipe_type"`
	SwipeDate    time.Time  `json:"swipe_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type Match struct {
	ID        uint      `json:"id"`
	User1ID   uint      `json:"user1_id"`
	User2ID   uint      `json:"user2_id"`
	MatchedAt time.Time `json:"matched_at"`
}

type MatchDetail struct {
	MatchID       uint                   `json:"match_id"`
	MatchedUserID uint                   `json:"matched_user_id"`
	Profile       *profileDomain.Profile `json:"profile"`
	MatchedAt     time.Time              `json:"matched_at"`
}

type SwipeRepository interface {
	CreateSwipe(ctx context.Context, swipe *Swipe) error
	CreateSwipeWithTx(ctx context.Context, tx database.DBTX, swipe *Swipe) error
	GetDailySwipeCount(ctx context.Context, userID uint) (int, error)
	GetDailySwipeCountWithTx(ctx context.Context, tx database.DBTX, userID uint) (int, error)
	HasSwiped(ctx context.Context, userID, targetUserID uint) (bool, error)
	HasSwipedWithTx(ctx context.Context, tx database.DBTX, userID, targetUserID uint) (bool, error)
	GetSwipe(ctx context.Context, userID, targetUserID uint) (*Swipe, error)
	GetSwipeWithTx(ctx context.Context, tx database.DBTX, userID, targetUserID uint) (*Swipe, error)
	CreateMatch(ctx context.Context, user1ID, user2ID uint) (*Match, error)
	CreateMatchWithTx(ctx context.Context, tx database.DBTX, user1ID, user2ID uint) (*Match, error)
	GetMatch(ctx context.Context, user1ID, user2ID uint) (*Match, error)
	GetMatchesByUserID(ctx context.Context, userID uint) ([]MatchDetail, error)
}

type SwipeResponse struct {
	Swipe   *Swipe `json:"swipe"`
	IsMatch bool   `json:"is_match"`
	Match   *Match `json:"match,omitempty"`
}

type SwipeService interface {
	SwipeUser(ctx context.Context, userID, targetUserID uint, swipeType SwipeType) (*SwipeResponse, error)
	GetMatches(ctx context.Context, userID uint) ([]MatchDetail, error)
}
