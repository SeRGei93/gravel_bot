package table

import (
	"time"
)

// Gift — подарок участнику на событии
type Gift struct {
	ID           uint `gorm:"primaryKey;autoIncrement"`
	UserID       int64
	EventID      uint
	Content      string `gorm:"type:text"`
	MediaGroupId string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	User         User
	Files        []File `gorm:"foreignKey:EntityId;references:ID"`
}
