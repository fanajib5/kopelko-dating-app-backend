package repositories

import (
	model "kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type PremiumFeatureRepository interface {
	GetFeatureByID(id int) (*model.PremiumFeature, error)
}

type premiumFeatureRepo struct {
	db *gorm.DB
}

func NewPremiumFeatureRepository(db *gorm.DB) *premiumFeatureRepo {
	return &premiumFeatureRepo{db: db}
}

func (repo *premiumFeatureRepo) GetFeatureByID(id int) (*model.PremiumFeature, error) {
	var feature model.PremiumFeature
	err := repo.db.First(&feature, id).Error
	return &feature, err
}
