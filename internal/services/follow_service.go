package services

import (
	"twitter-system-design/internal/repositories"
)

type FollowService struct {
	Repo *repositories.FollowRepository
}

func NewFollowService(repo *repositories.FollowRepository) *FollowService {
	return &FollowService{Repo: repo}
}

func (s *FollowService) FollowUser(
	followerID int,
	followeeID int,
) error {
	if followerID == followeeID {
		return nil
	}
	return s.Repo.FollowUser(followerID, followeeID)
}

func (s *FollowService) UnfollowUser(
	followerID int,
	followeeID int,
) error {
	if followerID == followeeID {
		return nil
	}
	return s.Repo.UnfollowUser(followerID, followeeID)
}

func (s *FollowService) GetFollowing(id int) ([]int, error) {
	if id <= 0 {
		return nil, nil
	}
	return s.Repo.GetFollowing(id)
}
