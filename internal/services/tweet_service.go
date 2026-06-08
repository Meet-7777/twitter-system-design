package services

import (
	"errors"
	"strings"
	"twitter-system-design/internal/models"
	"twitter-system-design/internal/repositories"
)

type TweetService struct {
	Repo *repositories.TweetRepository
}

func NewTweetService(repo *repositories.TweetRepository) *TweetService {
	return &TweetService{Repo: repo}
}

func (s *TweetService) CreateTweet(tweet *models.Tweet) error {
	if tweet.UserID <= 0 {
		return errors.New("invalid user id")
	}

	if strings.TrimSpace(tweet.Content) == "" {
		return errors.New("tweet cannot be empty")
	}
	return s.Repo.CreateTweet(tweet.UserID, tweet.Content)
}

func (s *TweetService) GetTweetsByUserID(userID int) ([]string, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.Repo.GetTweetsByUserID(userID)
}
