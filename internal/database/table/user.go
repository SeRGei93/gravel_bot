package table

import (
	"time"
)

type User struct {
	ID         int64  `gorm:"primaryKey;autoIncrement:false"`
	NickName   string `gorm:"default:''"`
	FirstName  string `gorm:"default:''"`
	LastName   string `gorm:"default:''"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserEvents []UserEvent
}
