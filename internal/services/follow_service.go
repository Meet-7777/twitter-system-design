package services

import "twitter-system-design/internal/repository"

type FollowService struct {
	followRepo *repository.FollowRepository
}

func NewFollowService(followRepo *repository.FollowRepository) *FollowService {
	return &FollowService{followRepo: followRepo}
}

func (s *FollowService) FollowUser(followerID int, followeeID int) error {
	return s.followRepo.Create(followerID, followeeID)
}
