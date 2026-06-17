package queue

type FanoutEvent struct {
	TweetID     int
	FollowerIDs []int
	Score       float64
}

var FanoutQueue = make(chan FanoutEvent, 1000)
