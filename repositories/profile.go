package repositories

import (
	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	BeginTx() *gorm.DB
	CreateProfileTx(tx *gorm.DB, profile *models.Profile) error
	FindByID(id uint) (*models.Profile, error)
	UpdateIsPremiumTx(tx *gorm.DB, userID uint, isPremium bool) error
	FindRandom() (*models.Profile, error)
}

type profileRepo struct {
	db *gorm.DB
}

func (r *profileRepo) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func NewProfileRepository(db *gorm.DB) *profileRepo {
	return &profileRepo{db: db}
}

func (r *profileRepo) CreateProfileTx(tx *gorm.DB, profile *models.Profile) error {
	return tx.Create(profile).Error
}

func (r *profileRepo) FindByID(id uint) (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.First(&profile, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateIsPremium sets the IsPremium field for a user's profile
func (r *profileRepo) UpdateIsPremiumTx(tx *gorm.DB, userID uint, isPremium bool) error {
	return tx.Model(&models.Profile{}).Where("user_id = ?", userID).
		Update("is_premium", isPremium).Error
}

// FindRandom returns a random profile
func (r *profileRepo) FindRandom() (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.Order("RANDOM()").First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
