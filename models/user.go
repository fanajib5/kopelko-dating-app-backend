package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string  `gorm:"unique;not null;column:email"`
	Password   string  `gorm:"not null;column:password_hash"`
	IsVerified bool    `gorm:"default:false;column:is_verified"`
	Profile    Profile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Specify table name if different from struct name
func (User) TableName() string {
	return "users"
}
