package repository

import (
	"gravel_bot/internal/database/table"

	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (r *UserRepository) Init() *UserRepository {
	r.database.AutoMigrate(&table.User{})
	return r
}

func (r *UserRepository) CreateUser(tu table.User) error {
	err := r.database.Create(&tu).Error
	return err
}

func (r *UserRepository) UpdateUser(tu table.User) error {
	err := r.database.Save(&tu).Error
	return err
}

func (r *UserRepository) FindUser(id int64) (table.User, error) {
	tu := table.User{ID: id}
	err := r.database.Where(tu).Take(&tu).Error
	if err != nil {
		return tu, err
	}

	return tu, nil
}

func (r *UserRepository) DeleteUser(id int64) error {
	tu := table.User{ID: id}
	return r.database.Delete(tu).Error
}

func (r *UserRepository) GetAllUsers() ([]table.User, error) {
	var users []table.User
	err := r.database.Find(&users).Error
	return users, err
}
