package repositories

import (
	"time"

	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type ProfileViewRepository interface {
	CreateProfileView(view *models.ProfileView) error
	GetUnviewedProfiles(userID uint, limit int) ([]models.Profile, error)
}

type profileViewRepo struct {
	db *gorm.DB
}

func NewProfileViewRepository(db *gorm.DB) *profileViewRepo {
	return &profileViewRepo{db: db}
}

// Fetches up to 10 profiles that haven't been viewed by the user today.
func (r *profileViewRepo) GetUnviewedProfiles(userID uint, limit int) ([]models.Profile, error) {
	var profiles []models.Profile
	today := time.Now().Format("2006-01-02")

	if err := r.db.Raw(`
        SELECT p.* FROM profiles p
        LEFT JOIN profile_views pv ON p.user_id = pv.viewed_user_id 
            AND pv.user_id = ? AND pv.view_date = ?
        WHERE pv.id IS NULL
        ORDER BY RANDOM()
        LIMIT ?`, userID, today, limit).Scan(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

// CreateProfileView creates a new profile view in the database
func (r *profileViewRepo) CreateProfileView(view *models.ProfileView) error {
	return r.db.Create(view).Error
}
