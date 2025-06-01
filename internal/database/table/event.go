package table

import (
	"time"
)

// Event — велогонка/мероприятие
type Event struct {
	ID        uint `gorm:"primaryKey"`
	Active    bool `gorm:"default:true"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Gifts     []Gift
}
