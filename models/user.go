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
	parts := strings.Split(u.Email, "@")
	if len(parts) != 2 {
		// return the original email if it doesn't have exactly one '@' character
		return u.Email
	}
	maskedLocal := maskLocalPart(parts[0])
	return maskedLocal + "@*****"
}

// maskLocalPart masks all but the first three characters of the local part
func maskLocalPart(local string) string {
	if len(local) <= 3 {
		return local
	}
	return local[:3] + strings.Repeat("*", len(local)-3)
}
