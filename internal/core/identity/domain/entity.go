package domain

import (
	"context"
	"time"

	"kopelko-dating-app-backend/internal/platform/database"
)

// User represents the core account identity.
type User struct {
	ID           uint       `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	IsVerified   bool       `json:"is_verified"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// UserRegisteredPayload is the payload passed to the "user.registered" action hook.
type UserRegisteredPayload struct {
	User      *User
	Name      string
	Age       int
	Gender    string
	Location  string
	Interests []string
	Photos    []string
	Tx        database.DBTX
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	CreateWithTx(ctx context.Context, tx database.DBTX, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	SetVerified(ctx context.Context, id uint, isVerified bool) error
}

type IdentityService interface {
	Register(ctx context.Context, email, password, name string, age int, gender, location string, interests, photos []string) (*User, string, error)
	Login(ctx context.Context, email, password string) (*User, string, error)
}
