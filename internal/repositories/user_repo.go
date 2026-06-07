package repositories

import (
	"database/sql"
	"twitter-system-design/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.QueryRow(
		"INSERT INTO users(username) VALUES($1) RETURNING id",
		user.Username,
	).Scan(&user.ID)
}
