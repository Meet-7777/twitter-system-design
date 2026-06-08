package repositories

import "database/sql"

type TweetRepository struct {
	DB *sql.DB
}

func NewTweetRepositrory(db *sql.DB) *TweetRepository {
	return &TweetRepository{DB: db}
}

func (r *TweetRepository) CreateTweet(userID int, content string) error {
	_, err := r.DB.Exec(
		"INSERT INTO tweets(user_id, content) VALUES($1, $2)",
		userID, content,
	)
	return err
}

func (r *TweetRepository) GetTweetsByUserID(userID int) ([]string, error) {
	rows, err := r.DB.Query(
		"SELECT content FROM tweets WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		tweets = append(tweets, content)
	}
	return tweets, nil
}
