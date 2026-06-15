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

func NewTweetService(
	repo *repositories.TweetRepository,
) *TweetService {

	return &TweetService{
		Repo: repo,
	}

}

func (s *TweetService) CreateTweet(
	tweet *models.Tweet,
) (int, error) {

	if tweet.UserID <= 0 {
		return 0, errors.New(
			"invalid user id",
		)
	}

	if strings.TrimSpace(tweet.Content) == "" {
		return 0, errors.New(
			"tweet empty",
		)
	}

	if len(tweet.Content) > 280 {
		return 0, errors.New(
			"tweet too long",
		)
	}

	return s.Repo.CreateTweet(
		tweet.UserID,
		tweet.Content,
	)

}
