package repositories

import (
	model "kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	BeginTx() *gorm.DB
	CreateUserTx(tx *gorm.DB, user *model.User) error
	FindByEmail(email string) (*model.User, error)
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

func (r *userRepo) CreateUserTx(tx *gorm.DB, user *model.User) error {
	return tx.Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
