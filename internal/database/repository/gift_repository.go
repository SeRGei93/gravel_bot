package repository

import (
	"gravel_bot/internal/database/table"

	"time"

	"gorm.io/gorm"
)

type GiftRepository struct {
	database *gorm.DB
}

type GiftDto struct {
	ID        int64
	NickName  string
	FirstName string
	LastName  string
	Content   string
	CreatedAt time.Time
	Files     []table.File
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

func (r *GiftRepository) FindGiftByMediaGroup(id string) (table.Gift, error) {
	var row table.Gift
	err := r.database.Where("media_group_id = ?", id).First(&row).Error
	return row, err
}

func (r *GiftRepository) FindGiftsByEvent(id uint) ([]table.Gift, error) {
	var rows []table.Gift
	err := r.database.Where("event_id = ?", id).Preload("User").Preload("Files").Order("created_at DESC").Find(&rows).Error
	return rows, err
}

// ExportGifts выгружает подарки в CSV-файл
func (r *GiftRepository) ExportGifts(eventID uint) ([]GiftDto, error) {
	var results []GiftDto

	gifts, err := r.FindGiftsByEvent(eventID)
	if err != nil {
		return nil, err
	}

	for _, gift := range gifts {

		dto := GiftDto{
			ID:        gift.User.ID,
			NickName:  gift.User.NickName,
			FirstName: gift.User.FirstName,
			LastName:  gift.User.LastName,
			Content:   gift.Content,
			CreatedAt: gift.CreatedAt,
			Files:     gift.Files,
		}

		results = append(results, dto)
	}

	return results, nil
}
