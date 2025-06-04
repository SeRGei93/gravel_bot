package repository

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gravel_bot/internal/database/table"
	"os"
	"time"

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

// ExportEventParticipantsCSV выгружает участников события в CSV-файл
func (r *UserEventRepository) ExportEventParticipantsCSV(eventID uint, outputPath string) error {
	type row struct {
		NickName  string
		FirstName string
		LastName  string
		CreatedAt time.Time
		Bike      string
		Gift      string
		Result    string
	}

	var results []row

	err := r.database.Raw(`
		SELECT 
			u.nick_name,
			u.first_name,
			u.last_name,
			ue.created_at,
			ue.bike,
			ue.result_link AS result,
			CASE WHEN g.id IS NOT NULL THEN 'Y' ELSE '' END AS gift
		FROM user_events ue
		JOIN users u ON u.id = ue.user_id
		LEFT JOIN gifts g ON g.user_id = ue.user_id AND g.event_id = ue.event_id
		WHERE ue.event_id = ?
		GROUP BY ue.user_id
		ORDER BY ue.created_at
	`, eventID).Scan(&results).Error
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"username", "first_name", "last_name", "registered_at", "bike_type", "gift", "result"})

	for _, row := range results {
		writer.Write([]string{
			row.NickName,
			row.FirstName,
			row.LastName,
			row.CreatedAt.Format("02.01.2006 15:04:05"),
			row.Bike,
			row.Gift,
			row.Result,
		})
	}

	return nil
}
