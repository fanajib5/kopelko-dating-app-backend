package repositories

import (
	model "kopelko-dating-app-backend/models"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscription(subscription *model.Subscription) error
	HasFeature(userID uint, featureName string) (bool, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *subscriptionRepo {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) CreateSubscription(subscription *model.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepo) HasFeature(userID uint, featureName string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Subscription{}).
		Joins("JOIN premium_features ON subscriptions.feature_id = premium_features.id").
		Where("subscriptions.user_id = ? AND premium_features.feature_name = ? AND subscriptions.end_date > ?", userID, featureName, time.Now()).
		Count(&count).Error

	return count > 0, err
}
