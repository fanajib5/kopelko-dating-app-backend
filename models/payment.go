package models

import "time"

type Payment struct {
	ID              uint      `gorm:"primaryKey;column:id"`
	UserID          uint      `gorm:"not null;column:user_id"`
	Amount          float64   `gorm:"not null;column:amount;check:amount > 0"`
	Currency        string    `gorm:"default:'USD';column:currency"`
	PaymentDate     time.Time `gorm:"default:current_timestamp;column:payment_date"`
	PaymentStatus   string    `gorm:"column:payment_status"`
	PaymentProvider string    `gorm:"column:payment_provider"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (Payment) TableName() string {
	return "payments"
}
