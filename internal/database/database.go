package database

import (
	"fmt"
	Repository "gravel_bot/internal/database/repository"

	"gorm.io/gorm"
)

type Database struct {
	User      Repository.UserRepository
	Event     Repository.EventRepository
	UserEvent Repository.UserEventRepository
	Gift      Repository.GiftRepository
	File      Repository.FileRepository
}

func InitDatabase(dialector gorm.Dialector) Database {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	return Database{
		User:      *Repository.NewUserRepository(db).Init(),
		Event:     *Repository.NewEventRepository(db).Init(),
		UserEvent: *Repository.NewUserEventRepository(db).Init(),
		Gift:      *Repository.NewGiftRepository(db).Init(),
		File:      *Repository.NewFileRepository(db).Init(),
	}
}
