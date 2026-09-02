package domain

import (
	"context"
	"time"

	"kopelko-dating-app-backend/internal/platform/database"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Profile struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	Name       string     `json:"name"`
	Age        int        `json:"age"`
	Bio        string     `json:"bio"`
	Gender     Gender     `json:"gender"`
	Location   string     `json:"location"`
	Interests  []string   `json:"interests"`
	Photos     []string   `json:"photos"`
	IsPremium  bool       `json:"is_premium"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type ProfileView struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ViewedUserID uint      `json:"viewed_user_id"`
	SwipeID      *uint     `json:"swipe_id,omitempty"`
	ViewDate     time.Time `json:"view_date"`
}

type DiscoveryFilter struct {
	Gender *string `json:"gender,omitempty"`
	MinAge *int    `json:"min_age,omitempty"`
	MaxAge *int    `json:"max_age,omitempty"`
	Limit  int     `json:"limit"`
}

type UpdateProfileRequest struct {
	Name      *string   `json:"name,omitempty" validate:"omitempty"`
	Age       *int      `json:"age,omitempty" validate:"omitempty,gte=18"`
	Bio       *string   `json:"bio,omitempty" validate:"omitempty"`
	Gender    *string   `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Location  *string   `json:"location,omitempty" validate:"omitempty"`
	Interests []string  `json:"interests,omitempty" validate:"omitempty"`
	Photos    []string  `json:"photos,omitempty" validate:"omitempty"`
}

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	CreateWithTx(ctx context.Context, tx database.DBTX, profile *Profile) error
	GetByUserID(ctx context.Context, userID uint) (*Profile, error)
	Update(ctx context.Context, profile *Profile) error
	RecordView(ctx context.Context, userID, viewedUserID uint, swipeID *uint) error
	RecordViewWithTx(ctx context.Context, tx database.DBTX, userID, viewedUserID uint, swipeID *uint) error
	GetDailyViewCount(ctx context.Context, userID uint) (int, error)
	GetRandomProfiles(ctx context.Context, currentUserID uint, filter DiscoveryFilter) ([]Profile, error)
}

type ProfileService interface {
	GetMyProfile(ctx context.Context, userID uint) (*Profile, error)
	GetRandomProfiles(ctx context.Context, currentUserID uint, filter DiscoveryFilter) ([]Profile, error)
	CreateProfile(ctx context.Context, profile *Profile) error
	UpdateMyProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*Profile, error)
}
