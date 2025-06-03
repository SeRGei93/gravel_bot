package repository

import (
	"gravel_bot/internal/database/table"

	"gorm.io/gorm"
)

type FileRepository struct {
	database *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{database: db}
}

func (r *FileRepository) Init() *FileRepository {
	r.database.AutoMigrate(&table.Gift{})
	return r
}

func (r *FileRepository) CreateFile(tu table.File) error {
	err := r.database.Create(&tu).Error
	return err
}
