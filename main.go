package main

import (
	"twitter-system-design/internal/database"
	"twitter-system-design/internal/handlers"
	"twitter-system-design/internal/kafka"
	"twitter-system-design/internal/repository"
	"twitter-system-design/internal/services"
	"twitter-system-design/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.NewPostgres()
	rdb := database.NewRedis()

	producer := kafka.NewProducer("localhost:9092", "tweets")
	consumer := kafka.NewConsumer("localhost:9092", "tweets", "fanout-group")
	defer producer.Close()
	defer consumer.Close()

	tweetRepo := repository.NewTweetRepository(db)
	followRepo := repository.NewFollowRepository(db)
	userRepo := repository.NewUserRepository(db)

	tweetService := services.NewTweetService(tweetRepo, followRepo, userRepo, rdb, producer)
	timelineService := services.NewTimelineService(tweetRepo, followRepo, rdb)
	userService := services.NewUserService(userRepo)
	followService := services.NewFollowService(followRepo)

	tweetHandler := handlers.NewTweetHandler(tweetService)
	timelineHandler := handlers.NewTimelineHandler(timelineService)
	userHandler := handlers.NewUserHandler(userService)
	followHandler := handlers.NewFollowHandler(followService)

	fanoutWorker := worker.NewFanoutWorker(rdb, consumer)
	go fanoutWorker.Start()

	r := gin.Default()
	r.POST("/tweet", tweetHandler.CreateTweet)
	r.GET("/feed/:id", timelineHandler.GetFeed)
	r.POST("/user", userHandler.CreateUser)
	r.POST("/follow", followHandler.FollowUser)
	r.Run(":8080")
}
