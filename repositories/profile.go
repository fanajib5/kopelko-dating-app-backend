package repositories

import (
	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfileTx(tx *gorm.DB, profile *models.Profile) error
	FindByID(id string) (*models.Profile, error)
}

type profileRepo struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *profileRepo {
	return &profileRepo{db: db}
}

func (r *profileRepo) CreateProfileTx(tx *gorm.DB, profile *models.Profile) error {
	return tx.Create(profile).Error
}

func (r *profileRepo) FindByID(id string) (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.First(&profile, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
