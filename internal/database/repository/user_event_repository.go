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
func (r *UserEventRepository) RegisterUserToEvent(userID int64, eventID uint, active bool, bike string) error {
	_, err := r.FindUserToEvent(userID, eventID)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return r.database.Create(&table.UserEvent{
			UserID:  userID,
			EventID: eventID,
			Active:  active,
			Bike:    bike,
		}).Error
	}

	return fmt.Errorf("пользователь уже зарегистрирован")
}

// Зарегистрировать пользователя на событие
func (r *UserEventRepository) UnRegisterUserToEvent(userID int64, eventID uint) error {
	existing, err := r.FindUserToEvent(userID, eventID)

	if err != nil {
		return err
	}

	return r.database.Delete(existing).Error
}

func (r *UserEventRepository) FindUserToEvent(userID int64, eventID uint) (*table.UserEvent, error) {
	var existing table.UserEvent
	err := r.database.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&existing).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &existing, nil
}
