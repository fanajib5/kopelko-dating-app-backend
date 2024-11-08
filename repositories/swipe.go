package repositories

import (
	"time"

	model "kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type SwipeRepository interface {
	CreateSwipe(swipe *model.Swipe) error
	GetDailySwipes(userID uint, date time.Time) (int64, error)
	HasSwipedToday(userID, targetUserID uint, date time.Time) (bool, error)
}

type swipeRepository struct {
	db *gorm.DB
}

// NewSwipeRepository initializes a new repository
func NewSwipeRepository(db *gorm.DB) *swipeRepository {
	return &swipeRepository{db: db}
}

// CreateSwipe creates a new swipe in the database
func (r *swipeRepository) CreateSwipe(swipe *model.Swipe) error {
	return r.db.Create(swipe).Error
}

// GetDailySwipes retrieves the count of swipes for a user on a given day
func (r *swipeRepository) GetDailySwipes(userID uint, date time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&model.Swipe{}).
		Where("user_id = ? AND swipe_date = ?", userID, date).
		Count(&count).Error
	return count, err
}

// HasSwipedToday checks if a user has swiped on a target user today
func (r *swipeRepository) HasSwipedToday(userID, targetUserID uint, date time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&model.Swipe{}).
		Where("user_id = ? AND target_user_id = ? AND swipe_date = ?", userID, targetUserID, date).
		Count(&count).Error
	return count > 0, err
}
