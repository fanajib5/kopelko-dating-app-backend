package repositories

import (
	"time"

	model "kopelko-dating-app-backend/models"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscriptionTx(tx *gorm.DB, subscription *model.Subscription) error
	HasFeature(userID uint, featureName string) (bool, error)
	GetActiveSubscription(userID uint) (bool, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *subscriptionRepo {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) CreateSubscriptionTx(tx *gorm.DB, subscription *model.Subscription) error {
	return tx.Create(subscription).Error
}

func (r *subscriptionRepo) HasFeature(userID uint, featureName string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Subscription{}).
		Joins("JOIN premium_features ON subscriptions.feature_id = premium_features.id").
		Where("subscriptions.user_id = ? AND premium_features.feature_name = ? AND subscriptions.end_date > ?", userID, featureName, time.Now()).
		Count(&count).Error

	return count > 0, err
}

// GetActiveSubscription checks if a user has an active subscription
func (r *subscriptionRepo) GetActiveSubscription(userID uint) (bool, error) {
	var count int64
	err := r.db.Where("user_id = ? AND end_date > ?", userID, true, time.Now()).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}
