package models

import (
	"time"

	"gorm.io/gorm"
)

type ProfileView struct {
	gorm.Model
	UserID       uint      `gorm:"not null"`
	ViewedUserID uint      `gorm:"not null"`
	SwipeID      uint      `gorm:"default:null"`
	ViewDate     time.Time `gorm:"default:CURRENT_DATE"`
}

// Struct for managing swipe and profile view data together
type SwipeAndViewData struct {
	Swipe       Swipe
	ProfileView ProfileView
}

func (ProfileView) TableName() string {
	return "profile_views"
}
