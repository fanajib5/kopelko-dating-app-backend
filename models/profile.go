package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID    uint     `gorm:"unique;not null;column:user_id"`
	Name      string   `gorm:"not null;column:name"`
	Age       int      `gorm:"column:age;check:age >= 18"`
	Bio       string   `gorm:"column:bio"`
	Gender    string   `gorm:"column:gender"`
	Location  string   `gorm:"column:location"`
	Interests string   `gorm:"column:interests"`
	Photos    []string `gorm:"type:text[];column:photos"`
	IsPremium bool     `gorm:"default:false;column:is_premium"`
}

func (Profile) TableName() string {
	return "profiles"
}
