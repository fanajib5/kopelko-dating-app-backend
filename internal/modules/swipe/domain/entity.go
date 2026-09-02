package domain

import (
	"context"
	"time"
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

type SwipeRepository interface {
	CreateSwipe(ctx context.Context, swipe *Swipe) error
	GetDailySwipeCount(ctx context.Context, userID uint) (int, error)
	HasSwiped(ctx context.Context, userID, targetUserID uint) (bool, error)
	GetSwipe(ctx context.Context, userID, targetUserID uint) (*Swipe, error)
	CreateMatch(ctx context.Context, user1ID, user2ID uint) (*Match, error)
	GetMatch(ctx context.Context, user1ID, user2ID uint) (*Match, error)
}

type SwipeResponse struct {
	Swipe   *Swipe `json:"swipe"`
	IsMatch bool   `json:"is_match"`
	Match   *Match `json:"match,omitempty"`
}

type SwipeService interface {
	SwipeUser(ctx context.Context, userID, targetUserID uint, swipeType SwipeType) (*SwipeResponse, error)
}
