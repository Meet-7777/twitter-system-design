package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repositories"

	"github.com/redis/go-redis/v9"
)

type TweetService struct {
	Repo       *repositories.TweetRepository
	FollowRepo *repositories.FollowRepository
	Cache      *redis.Client
}

func NewTweetService(
	repo *repositories.TweetRepository,
	followRepo *repositories.FollowRepository,
	cache *redis.Client,
) *TweetService {
	return &TweetService{
		Repo:       repo,
		FollowRepo: followRepo,
		Cache:      cache,
	}
}

func (s *TweetService) CreateTweet(
	tweet *models.Tweet,
) (int, error) {

	if tweet.UserID <= 0 {
		return 0, errors.New("invalid user id")
	}

	if strings.TrimSpace(tweet.Content) == "" {
		return 0, errors.New("tweet cannot be empty")
	}

	if len(tweet.Content) > 280 {
		return 0, errors.New("tweet exceeds 280 characters")
	}

	tweetID, err := s.Repo.CreateTweet(
		tweet.UserID,
		tweet.Content,
	)
	if err != nil {
		return 0, err
	}
	followers, err := s.FollowRepo.GetFollowers(
		tweet.UserID,
	)
	if err != nil {
		return tweetID, err
	}

	ctx := context.Background()
	for _, followerID := range followers {
		feedKey := fmt.Sprintf(
			"feed:%d",
			followerID,
		)
		err := s.Cache.LPush(
			ctx,
			feedKey,
			tweetID,
		).Err()
		if err != nil {
			fmt.Println(
				"redis feed push failed:",
				err,
			)

		}

	}
	return tweetID, nil
}

func (s *TweetService) GetTweetsByUserID(
	userID int,
) ([]models.Tweet, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	return s.Repo.GetTweetsByUserID(userID)
}
