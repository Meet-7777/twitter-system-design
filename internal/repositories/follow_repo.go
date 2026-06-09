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

func (r *FollowRepository) FollowUser(
	followerID,
	followeeID int,
) error {
	_, err := r.DB.Exec(
		`INSERT INTO follows(follower_id, followee_id)
		 VALUES($1, $2)`,
		followerID,
		followeeID,
	)

	return err
}

func (r *FollowRepository) UnfollowUser(
	followerID,
	followeeID int,
) error {
	_, err := r.DB.Exec(
		`DELETE FROM follows
		 WHERE follower_id = $1
		 AND followee_id = $2`,
		followerID,
		followeeID,
	)

	return err
}

func (r *FollowRepository) GetFollowing(
	userID int,
) ([]models.User, error) {

	rows, err := r.DB.Query(
		`
		SELECT u.id, u.username
		FROM users u
		INNER JOIN follows f
			ON u.id = f.followee_id
		WHERE f.follower_id = $1
		ORDER BY u.username
		`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err := rows.Scan(
			&user.ID,
			&user.Username,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
