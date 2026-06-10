package main

import (
	"log"
	"net/http"

	"twitter-system-design/internal/database"
	"twitter-system-design/internal/handlers"
	"twitter-system-design/internal/repositories"
	"twitter-system-design/internal/services"
)

func main() {
	redisClient := database.NewRedis()
	db := database.NewPostgres()
	defer redisClient.Close()
	defer db.Close()

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	followRepo := repositories.NewFollowRepository(db)
	tweetRepo := repositories.NewTweetRepository(db)
	timelineRepo := repositories.NewTimelineRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	followService := services.NewFollowService(followRepo)
	tweetService := services.NewTweetService(tweetRepo, followRepo, redisClient)
	timelineService := services.NewTimelineService(timelineRepo, redisClient)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	followHandler := handlers.NewFollowHandler(followService)
	tweetHandler := handlers.NewTweetHandler(tweetService)
	timelineHandler := handlers.NewTimelineHandler(timelineService)

	// User routes
	http.HandleFunc("/users", userHandler.CreateUser)

	// Follow routes
	http.HandleFunc("/follow", followHandler.FollowUser)
	http.HandleFunc("/following", followHandler.GetFollowing)

	// Tweet routes
	http.HandleFunc("/tweets", tweetHandler.Tweets)

	// Timeline route
	http.HandleFunc("/timeline", timelineHandler.GetTimeline)

	log.Println("🚀 Server running on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf(
			"Could not start server: %v",
			err,
		)
	}
}
