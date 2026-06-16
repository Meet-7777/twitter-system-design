package repository

import (
	"database/sql"
)

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) Create(followerId int, followeeId int) error {
	_, err := r.db.Exec(
		"INSERT INTO follows(follower_id, followee_id) VALUES ($1,$2)", followerId, followeeId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *FollowRepository) GetFollowerIDs(userID int) ([]int, error) {
	rows, err := r.db.Query(
		"SELECT follower_id FROM follows WHERE followee_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *FollowRepository) GetCelebrityFolloweeIDs(userID int) ([]int, error) {
	rows, err := r.db.Query(`
        SELECT f.followee_id FROM follows f
        JOIN users u ON u.id = f.followee_id
        WHERE f.follower_id = $1 AND u.is_celebrity = true`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}
