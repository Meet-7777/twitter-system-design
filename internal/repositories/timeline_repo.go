package repositories

import (
	"database/sql"
	"twitter-system-design/internal/models"
)

type TimelineRepository struct {
	DB *sql.DB
}

func NewTimelineRepository(db *sql.DB) *TimelineRepository {
	return &TimelineRepository{DB: db}
}

func (r *TimelineRepository) GetTimeline(userID int, limit int) ([]models.TimelineTweet, error) {
	rows, err := r.DB.Query(
		`SELECT t.id, u.username, t.content, t.created_at
		 FROM tweets t
		 JOIN users u ON u.id = t.user_id
		 WHERE t.user_id IN (
			SELECT followee_id FROM follows WHERE follower_id = $1
		 )
		 ORDER BY t.created_at DESC
		 LIMIT $2`,
		userID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.TimelineTweet

	for rows.Next() {
		var t models.TimelineTweet
		_ = rows.Scan(&t.ID, &t.Username, &t.Content, &t.CreatedAt)
		tweets = append(tweets, t)
	}

	return tweets, nil
}
