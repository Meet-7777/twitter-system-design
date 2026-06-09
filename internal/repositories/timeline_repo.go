package repositories

import (
	"database/sql"
	"twitter-system-design/internal/models"
)

type TimelineRepository struct {
	DB *sql.DB
}

func NewTimelineRepository(db *sql.DB) *TimelineRepository {
	return &TimelineRepository{
		DB: db,
	}
}

func (r *TimelineRepository) GetTimeline(
	userID int,
	limit int,
) ([]models.TimelineTweet, error) {

	query := `
		SELECT
			tweets.id,
			users.username,
			tweets.content,
			tweets.created_at
		FROM tweets
		INNER JOIN users
			ON tweets.user_id = users.id
		WHERE tweets.user_id IN (
			SELECT followee_id
			FROM follows
			WHERE follower_id = $1
		)
		ORDER BY tweets.created_at DESC
		LIMIT $2
	`

	rows, err := r.DB.Query(
		query,
		userID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.TimelineTweet

	for rows.Next() {
		var tweet models.TimelineTweet

		err := rows.Scan(
			&tweet.ID,
			&tweet.Username,
			&tweet.Content,
			&tweet.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tweets, nil
}
