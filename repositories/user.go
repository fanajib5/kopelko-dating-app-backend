package repositories

import (
	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	BeginTx() *gorm.DB
	CreateUserTx(tx *gorm.DB, user *models.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *userRepo) CreateUserTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}
