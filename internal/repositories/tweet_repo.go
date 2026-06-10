package repositories

import (
	"database/sql"
	"twitter-system-design/internal/models"
)

type TweetRepository struct {
	DB *sql.DB
}

func NewTweetRepository(db *sql.DB) *TweetRepository {
	return &TweetRepository{DB: db}
}

func (r *TweetRepository) CreateTweet(userID int, content string) (int, error) {
	var tweetID int
	err := r.DB.QueryRow(
		`
				INSERT INTO tweets(user_id, content)
				VALUES($1, $2)
				RETURNING id
				`,
		userID,
		content,
	).Scan(&tweetID)
	if err != nil {
		return 0, err
	}
	return tweetID, nil

}

func (r *TweetRepository) GetTweetsByUserID(userID int) ([]models.Tweet, error) {
	rows, err := r.DB.Query(
		"SELECT id, user_id, content, created_at FROM tweets WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tweets, nil
}
