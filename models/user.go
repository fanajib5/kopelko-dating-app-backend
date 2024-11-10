package models

import (
	"time"

	"kopelko-dating-app-backend/utils"

	"gorm.io/gorm"
)

type User struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Profile    Profile        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Email      string         `gorm:"unique;not null;column:email"`
	Password   string         `gorm:"not null;column:password_hash"`
	IsVerified bool           `gorm:"default:false;column:is_verified"`
	Token      string         `gorm:"-"`
}

// Specify table name if different from struct name
func (User) TableName() string {
	return "users"
}

// MaskEmail returns the email with the first three characters of the local part revealed and the rest masked, and the domain part fully masked
func (u *User) MaskEmail() string {
	if u.Email == "" {
		return ""
	}

	return utils.MaskEmail(u.Email)
}
