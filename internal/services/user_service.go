package services

import (
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(username string) (*models.User, error) {
	return s.userRepo.Create(username)
}
