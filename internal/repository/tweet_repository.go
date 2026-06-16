package repository

import (
	"database/sql"
	"twitter-system-design/internal/models"

	"github.com/lib/pq"
)

type TweetRepository struct {
	db *sql.DB
}

func NewTweetRepository(db *sql.DB) *TweetRepository {
	return &TweetRepository{db: db}
}

func (r *TweetRepository) Create(userID int, content string) (*models.Tweet, error) {
	tweet := &models.Tweet{}
	err := r.db.QueryRow("INSERT INTO tweets(user_id, content) VALUES($1,$2) RETURNING id, user_id, content, created_at", userID, content).Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tweet, nil
}

func (r *TweetRepository) GetByIDs(ids []int) ([]*models.TimelineTweet, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := r.db.Query(`
        SELECT t.id, u.username, t.content, t.created_at
        FROM tweets t
        JOIN users u ON u.id = t.user_id
        WHERE t.id = ANY($1)
        ORDER BY t.created_at DESC`,
		pq.Array(ids),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*models.TimelineTweet
	for rows.Next() {
		t := &models.TimelineTweet{}
		rows.Scan(&t.ID, &t.Username, &t.Content, &t.CreatedAt)
		tweets = append(tweets, t)
	}
	return tweets, nil
}

func (r *TweetRepository) GetByUserIDs(userIDs []int, limit int) ([]*models.TimelineTweet, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	rows, err := r.db.Query(`
        SELECT t.id, u.username, t.content, t.created_at
        FROM tweets t
        JOIN users u ON u.id = t.user_id
        WHERE t.user_id = ANY($1)
        ORDER BY t.created_at DESC LIMIT $2`,
		pq.Array(userIDs), limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*models.TimelineTweet
	for rows.Next() {
		t := &models.TimelineTweet{}
		rows.Scan(&t.ID, &t.Username, &t.Content, &t.CreatedAt)
		tweets = append(tweets, t)
	}
	return tweets, nil
}
