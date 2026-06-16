package models

import "time"

type TimelineTweet struct {
	ID        int
	Username  string
	Content   string
	CreatedAt time.Time
}
