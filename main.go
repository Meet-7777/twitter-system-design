package main

import (
	"twitter-system-design/internal/database"
	"twitter-system-design/internal/handlers"
	"twitter-system-design/internal/kafka"
	"twitter-system-design/internal/middleware"
	"twitter-system-design/internal/repository"
	"twitter-system-design/internal/services"
	"twitter-system-design/internal/token"
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

	tokenStore := token.NewTokenStore(rdb)

	tweetService := services.NewTweetService(tweetRepo, followRepo, userRepo, rdb, producer)
	timelineService := services.NewTimelineService(tweetRepo, followRepo, rdb)
	followService := services.NewFollowService(followRepo)
	authService := services.NewAuthService(userRepo, tokenStore)

	tweetHandler := handlers.NewTweetHandler(tweetService)
	timelineHandler := handlers.NewTimelineHandler(timelineService)
	followHandler := handlers.NewFollowHandler(followService)
	authHandler := handlers.NewAuthHandler(authService)

	fanoutWorker := worker.NewFanoutWorker(rdb, consumer)
	go fanoutWorker.Start()

	r := gin.Default()

	r.POST("/signup", authHandler.Signup)
	r.GET("/verify", authHandler.VerifyEmail)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.Refresh)
	r.POST("/logout", authHandler.Logout)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/tweet", tweetHandler.CreateTweet)
		protected.GET("/feed/:id", timelineHandler.GetFeed)
		protected.POST("/follow", followHandler.FollowUser)
	}

	r.Run(":8080")
}
