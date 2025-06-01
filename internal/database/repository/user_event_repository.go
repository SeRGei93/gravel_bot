package repository

import (
	"errors"
	"fmt"
	"gravel_bot/internal/database/table"

	"gorm.io/gorm"
)

type UserEventRepository struct {
	database *gorm.DB
}

func NewUserEventRepository(db *gorm.DB) *UserEventRepository {
	return &UserEventRepository{database: db}
}

func (r *UserEventRepository) Init() *UserEventRepository {
	r.database.AutoMigrate(&table.UserEvent{})
	return r
}

// Зарегистрировать пользователя на событие
func (r *UserEventRepository) RegisterUserToEvent(userID int64, eventID uint, active bool) error {
	var existing table.UserEvent
	err := r.database.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&existing).Error

	if err == nil {
		return fmt.Errorf("пользователь уже зарегистрирован")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.database.Create(&table.UserEvent{
		UserID:  userID,
		EventID: eventID,
		Active:  active,
	}).Error
}

// Зарегистрировать пользователя на событие
func (r *UserEventRepository) UnRegisterUserToEvent(userID int64, eventID uint) error {
	var existing table.UserEvent
	err := r.database.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.database.Delete(existing).Error
}
