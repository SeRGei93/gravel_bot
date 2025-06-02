package table

import (
	"database/sql"
	"time"
)

type UserEvent struct {
	UserID     int64 `gorm:"primaryKey"`
	EventID    uint  `gorm:"primaryKey"`
	Active     bool  `gorm:"default:true"`
	Bike       string
	ResultLink sql.NullString
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
