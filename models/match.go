package models

import "time"

type Match struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	User1ID   uint      `gorm:"not null;column:user1_id"`
	User2ID   uint      `gorm:"not null;column:user2_id"`
	MatchedAt time.Time `gorm:"default:current_timestamp;column:matched_at"`
	User1     User      `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE"`
	User2     User      `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE"`
}

func (Match) TableName() string {
	return "matches"
}
