package services

import (
	"errors"

	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repositories"
)

type TimelineService struct {
	Repo *repositories.TimelineRepository
}

func NewTimelineService(
	repo *repositories.TimelineRepository,
) *TimelineService {

	return &TimelineService{
		Repo: repo,
	}

}

func (s *TimelineService) GetTimeline(
	userID int,
	limit int,
) ([]models.TimelineTweet, error) {

	if userID <= 0 {
		return nil, errors.New(
			"invalid user id",
		)
	}

	if limit <= 0 {
		limit = 20
	}

	return s.Repo.GetTimeline(
		userID,
		limit,
	)

}
