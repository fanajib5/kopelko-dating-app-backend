package models

import "time"

type Subscription struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	UserID    uint      `gorm:"unique;not null;column:user_id"`
	FeatureID uint      `gorm:"not null;column:feature_id"`
	StartDate time.Time `gorm:"not null;column:start_date"`
	EndDate   time.Time `gorm:"not null;column:end_date"`
	AutoRenew bool      `gorm:"default:true;column:auto_renew"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
