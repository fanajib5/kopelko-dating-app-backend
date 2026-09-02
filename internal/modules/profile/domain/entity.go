package domain

import (
	"context"
	"time"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Profile struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	Name      string     `json:"name"`
	Age       int        `json:"age"`
	Bio       string     `json:"bio"`
	Gender    Gender     `json:"gender"`
	Location  string     `json:"location"`
	Interests []string   `json:"interests"`
	Photos    []string   `json:"photos"`
	IsPremium bool       `json:"is_premium"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ProfileView struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ViewedUserID uint      `json:"viewed_user_id"`
	SwipeID      *uint     `json:"swipe_id,omitempty"`
	ViewDate     time.Time `json:"view_date"`
}

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	GetByUserID(ctx context.Context, userID uint) (*Profile, error)
	Update(ctx context.Context, profile *Profile) error
	RecordView(ctx context.Context, userID, viewedUserID uint, swipeID *uint) error
	GetDailyViewCount(ctx context.Context, userID uint) (int, error)
	GetRandomProfiles(ctx context.Context, currentUserID uint, limit int) ([]Profile, error)
}

type ProfileService interface {
	GetMyProfile(ctx context.Context, userID uint) (*Profile, error)
	GetRandomProfiles(ctx context.Context, currentUserID uint) ([]Profile, error)
	CreateProfile(ctx context.Context, profile *Profile) error
}
