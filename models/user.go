package models

import (
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Profile    Profile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Email      string  `gorm:"unique;not null;column:email"`
	Password   string  `gorm:"not null;column:password_hash"`
	IsVerified bool    `gorm:"default:false;column:is_verified"`
	Token      string  `gorm:"-"`
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

	parts := strings.Split(u.Email, "@")
	if len(parts) != 2 {
		// return the original email if it doesn't have exactly one '@' character
		return u.Email
	}
	local := parts[0]
	if len(local) <= 3 {
		return local + "@*****"
	}

	maskedLocal := local[:3] + strings.Repeat("*", 5)
	return maskedLocal + "@*****"
}
