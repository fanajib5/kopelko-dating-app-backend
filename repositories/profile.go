package repositories

import (
	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfileTx(tx *gorm.DB, profile *models.Profile) error
}

type profileRepo struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *profileRepo {
	return &profileRepo{DB: db}
}

func (r *profileRepo) CreateProfileTx(tx *gorm.DB, profile *models.Profile) error {
	return tx.Create(profile).Error
}
