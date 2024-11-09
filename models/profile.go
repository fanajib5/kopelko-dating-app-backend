package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID            uint           `json:"-" gorm:"primaryKey"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	UserID        uint           `json:"-" gorm:"unique;not null;column:user_id"`
	Name          string         `json:"name" gorm:"not null;column:name"`
	Age           int            `json:"age" gorm:"column:age;check:age >= 18"`
	Bio           string         `json:"bio" gorm:"column:bio"`
	Gender        string         `json:"gender" gorm:"column:gender"`
	Location      string         `json:"location" gorm:"column:location"`
	Interests     []string       `json:"interests" gorm:"column:interests"`
	Photos        []string       `json:"photos" gorm:"type:text[];column:photos"`
	IsPremium     bool           `json:"is_premium" gorm:"default:false;column:is_premium"`
	VerifiedLabel bool           `json:"verified_label"`
}

func (Profile) TableName() string {
	return "profiles"
}
