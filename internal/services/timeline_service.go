package services

import (
	"context"
	"fmt"
	"strconv"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repository"

	"github.com/redis/go-redis/v9"
)

type TimelineService struct {
	tweetRepo  *repository.TweetRepository
	followRepo *repository.FollowRepository
	rdb        *redis.Client
}

func NewTimelineService(tweetRepo *repository.TweetRepository, followRepo *repository.FollowRepository, rdb *redis.Client) *TimelineService {
	return &TimelineService{
		tweetRepo:  tweetRepo,
		followRepo: followRepo,
		rdb:        rdb,
	}

}

func (s *TimelineService) GetFeed(userID int, limit int) ([]*models.TimelineTweet, error) {
	ctx := context.Background()
	key := fmt.Sprintf("feed:%d", userID)
	results, err := s.rdb.ZRevRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	var tweetIDs []int

	for _, r := range results {
		id, _ := strconv.Atoi(r)
		tweetIDs = append(tweetIDs, id)
	}
	regularTweets, err := s.tweetRepo.GetByIDs(tweetIDs)
	if err != nil {
		return nil, err
	}
	celebsIds, err := s.followRepo.GetCelebrityFolloweeIDs(userID)
	if err != nil {
		return nil, err
	}
	celebsTweets, err := s.tweetRepo.GetByUserIDs(celebsIds, limit)
	if err != nil {
		return nil, err
	}
	return merge(regularTweets, celebsTweets, limit), nil
}

func merge(a, b []*models.TimelineTweet, limit int) []*models.TimelineTweet {
	result := make([]*models.TimelineTweet, 0, limit)
	i, j := 0, 0

	for i < len(a) && j < len(b) && len(result) < limit {
		if a[i].CreatedAt.After(b[j].CreatedAt) {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	for i < len(a) && len(result) < limit {
		result = append(result, a[i])
		i++
	}

	for j < len(b) && len(result) < limit {
		result = append(result, b[j])
		j++
	}

	return result
}
