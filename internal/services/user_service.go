package services

import (
	"errors"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {

	if user.Username == "" {
		return errors.New("username required")
	}

	return s.Repo.CreateUser(user)
}
