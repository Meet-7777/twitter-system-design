package services

import (
	"errors"
	"strings"
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
	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is required")
	}
	return s.Repo.CreateUser(user)
}