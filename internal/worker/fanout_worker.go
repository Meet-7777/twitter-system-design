package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"twitter-system-design/internal/kafka"
	"twitter-system-design/internal/queue"

	"github.com/redis/go-redis/v9"
)

type FanoutWorker struct {
	rdb      *redis.Client
	consumer *kafka.Consumer
}

func NewFanoutWorker(rdb *redis.Client, consumer *kafka.Consumer) *FanoutWorker {
	return &FanoutWorker{
		rdb:      rdb,
		consumer: consumer,
	}
}

func (w *FanoutWorker) Start() {
	ctx := context.Background()
	w.consumer.Consume(ctx, func(data []byte) error {
		var event queue.FanoutEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		return w.process(event)
	})
}

func (w *FanoutWorker) process(event queue.FanoutEvent) error {
	ctx := context.Background()
	member := strconv.Itoa(event.TweetID)

	for _, followerID := range event.FollowerIDs {
		key := fmt.Sprintf("feed:%d", followerID)
		err := w.rdb.ZAdd(ctx, key, redis.Z{
			Score:  event.Score,
			Member: member,
		}).Err()
		if err != nil {
			log.Printf("redis zadd error for follower %d: %v", followerID, err)
		}
	}
	return nil

}
