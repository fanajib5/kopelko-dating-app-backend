package repositories

import (
	model "kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfileTx(tx *gorm.DB, profile *model.Profile) error
	FindByID(id string) (*model.Profile, error)
}

type profileRepo struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *profileRepo {
	return &profileRepo{db: db}
}

func (r *profileRepo) CreateProfileTx(tx *gorm.DB, profile *model.Profile) error {
	return tx.Create(profile).Error
}

func (r *profileRepo) FindByID(id string) (*model.Profile, error) {
	var profile model.Profile
	if err := r.db.First(&profile, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
