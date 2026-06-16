package services

import (
	"context"
	"fmt"
	"strconv"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repository"

	"github.com/redis/go-redis/v9"
)

type TweetService struct {
	tweetRepo  *repository.TweetRepository
	followRepo *repository.FollowRepository
	userRepo   *repository.UserRepository
	rdb        *redis.Client
}

func NewTweetService(t *repository.TweetRepository, f *repository.FollowRepository, u *repository.UserRepository, rdb *redis.Client) *TweetService {
	return &TweetService{
		tweetRepo:  t,
		followRepo: f,
		userRepo:   u,
		rdb:        rdb,
	}
}

func (s *TweetService) CreateTweet(userId int, content string) (*models.Tweet, error) {
	user, err := s.userRepo.GetByID(userId)
	if err != nil {
		return nil, err
	}
	tweet, err := s.tweetRepo.Create(userId, content)
	if err != nil {
		return nil, err
	}
	if user.IsCelebrity {
		return tweet, nil
	}
	followerIDs, err := s.followRepo.GetFollowerIDs(userId)
	if err != nil {
		return nil, err
	}
	score := float64(tweet.CreatedAt.Unix())
	return tweet, s.fanout(followerIDs, tweet.ID, score)

}

func (s *TweetService) fanout(followerIDs []int, tweetID int, score float64) error {
	ctx := context.Background()
	member := strconv.Itoa(tweetID)

	for _, followerID := range followerIDs {
		key := fmt.Sprintf("feed:%d", followerID)
		err := s.rdb.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: member,
		}).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
