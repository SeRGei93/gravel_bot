package repository

import (
	"gravel_bot/internal/database/table"

	"fmt"

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

func (r *UserRepository) GetAllParticipants(eventID uint) ([]table.User, error) {
	type row struct {
		ID        int64
		NickName  string
		FirstName string
		LastName  string
		Gift      string
		Result    string
	}

	var results []row
	var users []table.User

	err := r.database.Raw(`
		SELECT
			u.id,
			u.nick_name,
			u.first_name,
			u.last_name,
			ue.result_link AS result,
			CASE WHEN g.id IS NOT NULL THEN 'Y' ELSE '' END AS gift
		FROM user_events ue
			JOIN users u ON u.id = ue.user_id
			LEFT JOIN gifts g ON g.user_id = ue.user_id AND g.event_id = ue.event_id
		WHERE ue.event_id = ? AND ue.result_link NOT NULL AND gift = 'Y' OR gift = 'Y'
		GROUP BY ue.user_id;
	`, eventID).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}

	for _, res := range results {
		users = append(users, table.User{
			ID:        res.ID,
			NickName:  res.NickName,
			FirstName: res.FirstName,
			LastName:  res.LastName,
		})
	}

	return users, err
}
