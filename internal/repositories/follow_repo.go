package repositories

import (
	"database/sql"
	"twitter-system-design/internal/models"
)

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (r *FollowRepository) FollowUser(followerID, followeeID int) error {
	_, err := r.DB.Exec(
		`INSERT INTO follows(follower_id, followee_id)
		 VALUES($1, $2)`,
		followerID,
		followeeID,
	)
	return err
}

func (r *FollowRepository) GetFollowers(userID int) ([]int, error) {
	rows, err := r.DB.Query(
		`SELECT follower_id
		 FROM follows
		 WHERE followee_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		followers = append(followers, id)
	}

	return followers, nil
}

func (r *FollowRepository) GetFollowing(userID int) ([]models.User, error) {
	rows, err := r.DB.Query(
		`SELECT u.id, u.username
		 FROM users u
		 JOIN follows f ON f.followee_id = u.id
		 WHERE f.follower_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		_ = rows.Scan(&u.ID, &u.Username)
		users = append(users, u)
	}

	return users, nil
}

func (r *FollowRepository) UnfollowUser(followerID, followeeID int) error {
	_, err := r.DB.Exec(
		`DELETE FROM follows
		 WHERE follower_id = $1 AND followee_id = $2`,
		followerID,
		followeeID,
	)
	return err
}
