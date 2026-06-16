package repository

import (
	"database/sql"
	"twitter-system-design/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`INSERT INTO users (username) VALUES ($1) RETURNING id, username, is_celebrity`, username,
	).Scan(&user.ID, &user.Username, &user.IsCelebrity)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, username, is_celebrity FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.IsCelebrity)
	if err != nil {
		return nil, err
	}
	return user, nil

}
