package repositories

import "database/sql"

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (r *FollowRepository) FollowUser(followerID, followeeID int) error {
	_, err := r.DB.Exec(
		"INSERT INTO follows(follower_id, followee_id) VALUES($1, $2)",
		followerID, followeeID,
	)
	return err
}

func (r *FollowRepository) UnfollowUser(followerID, followeeID int) error {
	_, err := r.DB.Exec(
		"DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2",
		followerID, followeeID,
	)
	return err
}

func (r *FollowRepository) GetFollowing(userID int) ([]int, error) {
	rows, err := r.DB.Query("SELECT followee_id FROM follows WHERE follower_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followees []int
	for rows.Next() {
		var followeeID int
		if err := rows.Scan(&followeeID); err != nil {
			return nil, err
		}
		followees = append(followees, followeeID)
	}
	return followees, nil

}
