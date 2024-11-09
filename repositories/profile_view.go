package repositories

import (
	"time"

	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type ProfileViewRepository interface {
	CreateProfileView(view *models.ProfileView) error
	GetUnviewedProfiles(userID uint) (*models.Profile, error)
	CreateSwipeAndView(data models.SwipeAndViewData) error
}

type profileViewRepo struct {
	db *gorm.DB
}

func NewProfileViewRepository(db *gorm.DB) *profileViewRepo {
	return &profileViewRepo{db: db}
}

// CreateSwipeAndView handles inserting a swipe action and logging a profile view in one transaction
func (r *profileViewRepo) CreateSwipeAndView(data models.SwipeAndViewData) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert the swipe record
		if err := tx.Create(&data.Swipe).Error; err != nil {
			return err
		}

		// Use the generated swipe ID for profile view record
		data.ProfileView.SwipeID = data.Swipe.ID

		// Insert the profile view record
		if err := tx.Create(&data.ProfileView).Error; err != nil {
			return err
		}

		return nil
	})
}

// Fetches up to 10 profiles that haven't been viewed by the user today.
func (r *profileViewRepo) GetUnviewedProfiles(userID uint) (*models.Profile, error) {
	var profile *models.Profile
	today := time.Now()

	if err := r.db.Raw(`
        SELECT p.* FROM profiles p
        LEFT JOIN profile_views pv ON p.user_id = pv.viewed_user_id 
            AND pv.user_id = ? AND pv.view_date = ? 
        WHERE pv.id IS NULL AND p.user_id <> ? 
        ORDER BY RANDOM()
        LIMIT 1`, userID, today, userID).Scan(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

// CreateProfileView creates a new profile view in the database
func (r *profileViewRepo) CreateProfileView(view *models.ProfileView) error {
	return r.db.Create(view).Error
}
