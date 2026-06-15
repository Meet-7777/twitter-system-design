package repositories

import (
	"database/sql"
	"twitter-system-design/internal/models"

	"github.com/lib/pq"
)

type TweetRepository struct {
	DB *sql.DB
}

func NewTweetRepository(db *sql.DB) *TweetRepository {
	return &TweetRepository{DB: db}
}

func (r *TweetRepository) CreateTweet(userID int, content string) (int, error) {
	var id int

	err := r.DB.QueryRow(
		`INSERT INTO tweets(user_id, content)
		 VALUES($1, $2)
		 RETURNING id`,
		userID,
		content,
	).Scan(&id)

	return id, err
}

func (r *TweetRepository) GetTweetsByUserID(userID int) ([]models.Tweet, error) {
	rows, err := r.DB.Query(
		`SELECT id, user_id, content, created_at
		 FROM tweets
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet

	for rows.Next() {
		var t models.Tweet
		if err := rows.Scan(&t.ID, &t.UserID, &t.Content, &t.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, t)
	}

	return tweets, nil
}

// used by TimelineService
func (r *TweetRepository) GetTweetsByIDs(ids []string) ([]models.TimelineTweet, error) {
	if len(ids) == 0 {
		return []models.TimelineTweet{}, nil
	}

	query := `
		SELECT t.id, u.username, t.content, t.created_at
		FROM tweets t
		JOIN users u ON u.id = t.user_id
		WHERE t.id = ANY($1)
		ORDER BY t.created_at DESC
	`

	rows, err := r.DB.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.TimelineTweet

	for rows.Next() {
		var t models.TimelineTweet
		if err := rows.Scan(&t.ID, &t.Username, &t.Content, &t.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, t)
	}

	return tweets, nil
}
