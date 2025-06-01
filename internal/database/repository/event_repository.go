package repository

import (
	"gravel_bot/internal/database/table"

	"gorm.io/gorm"
)

type EventRepository struct {
	database *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{database: db}
}

func (r *EventRepository) Init() *EventRepository {
	r.database.AutoMigrate(&table.Event{})
	r.InsertDefaultEvents()
	return r
}

func (r *EventRepository) CreateEvent(event *table.Event) error {
	return r.database.Create(event).Error
}

func (r *EventRepository) UpdateEvent(event *table.Event) error {
	return r.database.Save(event).Error
}

func (r *EventRepository) FindEvent(id uint) (*table.Event, error) {
	var event table.Event
	err := r.database.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) FindEventByName(name string) (*table.Event, error) {
	var event table.Event
	err := r.database.Where("name = ?", name).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetAllEvents() ([]table.Event, error) {
	var events []table.Event
	err := r.database.Find(&events).Error
	return events, err
}

func (r *EventRepository) InsertDefaultEvents() error {
	defaultEvents := []table.Event{
		{Name: "kamni200", Active: true},
		{Name: "kamni300", Active: false},
	}

	for _, evt := range defaultEvents {
		_, err := r.FindEventByName(evt.Name)
		if err == gorm.ErrRecordNotFound {
			if err := r.database.Create(&evt).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}
