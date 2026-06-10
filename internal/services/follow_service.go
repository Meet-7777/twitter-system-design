package services

import (
	"errors"

	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repositories"
)

type FollowService struct {
	Repo *repositories.FollowRepository
}

func NewFollowService(
	repo *repositories.FollowRepository,
) *FollowService {
	return &FollowService{
		Repo: repo,
	}
}

func (s *FollowService) FollowUser(
	followerID int,
	followeeID int,
) error {

	if followerID <= 0 || followeeID <= 0 {
		return errors.New("invalid user id")
	}

	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}

	return s.Repo.FollowUser(
		followerID,
		followeeID,
	)
}

func (s *FollowService) UnfollowUser(
	followerID int,
	followeeID int,
) error {

	if followerID <= 0 || followeeID <= 0 {
		return errors.New("invalid user id")
	}

	if followerID == followeeID {
		return errors.New("cannot unfollow yourself")
	}

	return s.Repo.UnfollowUser(
		followerID,
		followeeID,
	)
}

func (s *FollowService) GetFollowing(
	id int,
) ([]models.User, error) {

	if id <= 0 {
		return nil, errors.New("invalid user id")
	}

	return s.Repo.GetFollowing(id)
}

func (s *FollowService) GetFollowers(userID int) ([]int, error) {

	if userID <= 0 {
		return nil, errors.New("user id invalid")
	}
	return s.Repo.GetFollowers(userID)
}
