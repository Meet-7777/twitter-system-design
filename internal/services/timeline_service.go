package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repositories"

	"github.com/redis/go-redis/v9"
)

type TimelineService struct {
	Repo  *repositories.TimelineRepository
	Cache *redis.Client
}

func NewTimelineService(
	repo *repositories.TimelineRepository,
	cache *redis.Client,
) *TimelineService {
	return &TimelineService{
		Repo:  repo,
		Cache: cache,
	}
}

func (s *TimelineService) GetTimeline(
	userID int,
	limit int,
) ([]models.TimelineTweet, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if limit <= 0 {
		limit = 20
	}
	ctx := context.Background()

	cacheKey := fmt.Sprintf(
		"timeline:%d",
		userID,
	)

	cached, err := s.Cache.Get(ctx, cacheKey).Result()

	if err == nil {
		var tweets []models.TimelineTweet

		err := json.Unmarshal([]byte(cached), &tweets)

		if err == nil {
			fmt.Println("CACHE HIT")
			return tweets, nil
		}
	}
	fmt.Println("CACHE MISS")

	data, err := s.Repo.GetTimeline(
		userID,
		limit,
	)
	if err == nil {
		jsonData, err := json.Marshal(data)
		if err == nil {
			s.Cache.Set(ctx, cacheKey, jsonData, time.Minute*5)
		}

	}
	return data, nil
}
