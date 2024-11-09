package models

const (
	FeatureNameNoSwipeQuota  = "no_swipe_quota"
	FeatureNameVerifiedLabel = "verified_label"
)

type PremiumFeature struct {
	ID          uint   `gorm:"primaryKey;column:id"`
	FeatureName string `gorm:"not null;column:feature_name"`
	Description string `gorm:"not null;column:description"`
}

func (PremiumFeature) TableName() string {
	return "premium_features"
}
