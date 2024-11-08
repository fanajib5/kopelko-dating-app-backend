package models

import "time"

type Swipe struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	UserID       uint      `gorm:"not null;column:user_id"`
	TargetUserID uint      `gorm:"not null;column:target_user_id"`
	SwipeType    string    `gorm:"column:swipe_type;check:swipe_type IN ('left', 'right')"`
	SwipeDate    time.Time `gorm:"default:current_date;column:swipe_date"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	TargetUser   User      `gorm:"foreignKey:TargetUserID;constraint:OnDelete:CASCADE"`
}

func (Swipe) TableName() string {
	return "swipes"
}
