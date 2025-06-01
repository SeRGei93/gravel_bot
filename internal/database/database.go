package database

import (
	Repository "gravel_bot/internal/database/repository"

	"gorm.io/gorm"
)

type Database struct {
	User      Repository.UserRepository
	Event     Repository.EventRepository
	UserEvent Repository.UserEventRepository
}

func InitDatabase(dialector gorm.Dialector) Database {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return Database{
		User:      *Repository.NewUserRepository(db).Init(),
		Event:     *Repository.NewEventRepository(db).Init(),
		UserEvent: *Repository.NewUserEventRepository(db).Init(),
	}
}
