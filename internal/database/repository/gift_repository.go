package repository

import (
	"gravel_bot/internal/database/table"

	"gorm.io/gorm"
)

type GiftRepository struct {
	database *gorm.DB
}

func NewGiftRepository(db *gorm.DB) *GiftRepository {
	return &GiftRepository{database: db}
}

func (r *GiftRepository) Init() *GiftRepository {
	r.database.AutoMigrate(&table.Gift{})
	return r
}

func (r *GiftRepository) CreateGift(tu table.Gift) error {
	err := r.database.Create(&tu).Error
	return err
}
