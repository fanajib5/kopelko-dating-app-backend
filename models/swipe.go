package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	SwipeTypePass = "pass"
	SwipeTypeLike = "like"
)

type Swipe struct {
	ID           uint           `json:"-" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	UserID       uint           `gorm:"not null;column:user_id"`
	TargetUserID uint           `gorm:"not null;column:target_user_id"`
	SwipeType    string         `gorm:"column:swipe_type;check:swipe_type IN ('pass', 'like')"`
	SwipeDate    time.Time      `gorm:"default:current_date;column:swipe_date"`
	User         User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	TargetUser   User           `gorm:"foreignKey:TargetUserID;constraint:OnDelete:CASCADE"`
}

func (Swipe) TableName() string {
	return "swipes"
}
