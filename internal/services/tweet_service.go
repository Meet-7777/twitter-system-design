package services

import (
	"context"
	"log"
	"twitter-system-design/internal/kafka"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/queue"
	"twitter-system-design/internal/repository"

	"github.com/redis/go-redis/v9"
)

type TweetService struct {
	tweetRepo  *repository.TweetRepository
	followRepo *repository.FollowRepository
	userRepo   *repository.UserRepository
	rdb        *redis.Client
	producer   *kafka.Producer
}

func NewTweetService(t *repository.TweetRepository, f *repository.FollowRepository, u *repository.UserRepository, rdb *redis.Client, producer *kafka.Producer) *TweetService {
	return &TweetService{
		tweetRepo:  t,
		followRepo: f,
		userRepo:   u,
		rdb:        rdb,
		producer:   producer,
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
	event := queue.FanoutEvent{
		TweetID:     tweet.ID,
		FollowerIDs: followerIDs,
		Score:       float64(tweet.CreatedAt.Unix()),
	}
	if err := s.producer.Publish(context.Background(), event); err != nil {
		log.Printf("Kafka publish failed")
	} else {
		log.Printf("kafka publish successful")
	}
	return tweet, nil
}
