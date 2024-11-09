package repositories

import (
	"time"

	"kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscriptionTx(tx *gorm.DB, subscription *models.Subscription) error
	HasFeature(userID uint, featureName string) (bool, error)
	GetActiveSubscription(userID uint) (*models.Subscription, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *subscriptionRepo {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) CreateSubscriptionTx(tx *gorm.DB, subscription *models.Subscription) error {
	return tx.Create(subscription).Error
}

func (r *subscriptionRepo) HasFeature(userID uint, featureName string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Subscription{}).
		Joins("JOIN premium_features ON subscriptions.feature_id = premium_features.id").
		Where("subscriptions.user_id = ? AND premium_features.feature_name = ? AND subscriptions.end_date > ?", userID, featureName, time.Now()).
		Count(&count).Error

	return count > 0, err
}

// GetActiveSubscription checks if a user has an active subscription
func (r *subscriptionRepo) GetActiveSubscription(userID uint) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.Where("user_id = ? AND end_date > ?", userID, time.Now()).First(&subscription).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}
