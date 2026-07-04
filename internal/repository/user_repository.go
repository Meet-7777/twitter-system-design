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

func (r *UserRepository) Create(username, email, passwordHash string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, is_verified, is_celebrity`,
		username, email, passwordHash,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsVerified, &user.IsCelebrity)
	return user, err
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, username, email, is_verified, is_celebrity
		FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsVerified, &user.IsCelebrity)
	return user, err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, username, email, password_hash, is_verified, is_celebrity
		FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsVerified, &user.IsCelebrity)
	return user, err
}

func (r *UserRepository) UpdateVerified(id int) error {
	_, err := r.db.Exec(
		"UPDATE users SET is_verified = true WHERE id = $1", id,
	)
	return err
}
