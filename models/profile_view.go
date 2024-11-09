package models

import (
	"time"

	"gorm.io/gorm"
)

type ProfileView struct {
	gorm.Model
	UserID       uint      `gorm:"not null"`
	ViewedUserID uint      `gorm:"not null"`
	ViewDate     time.Time `gorm:"default:CURRENT_DATE"`
}

func (ProfileView) TableName() string {
	return "profile_views"
}
