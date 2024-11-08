package repositories

import (
	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

func (r *ProfileRepository) CreateProfile(profile *models.Profile) error {
	return r.DB.Create(profile).Error
}

func (r *ProfileRepository) CreateProfileTx(tx *gorm.DB, profile *models.Profile) error {
	return tx.Create(profile).Error
}
