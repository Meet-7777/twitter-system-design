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

	// repositories

	userRepo := repositories.NewUserRepository(db)

	followRepo := repositories.NewFollowRepository(db)

	tweetRepo := repositories.NewTweetRepository(db)

	timelineRepo := repositories.NewTimelineRepository(db)

	// services

	userService := services.NewUserService(userRepo)

	followService := services.NewFollowService(followRepo)

	tweetService := services.NewTweetService(tweetRepo)

	timelineService := services.NewTimelineService(
		timelineRepo,
	)

	// handlers

	userHandler := handlers.NewUserHandler(
		userService,
	)

	followHandler := handlers.NewFollowHandler(
		followService,
	)

	tweetHandler := handlers.NewTweetHandler(
		tweetService,
	)

	timelineHandler := handlers.NewTimelineHandler(
		timelineService,
	)

	http.HandleFunc(
		"/users",
		userHandler.CreateUser,
	)

	http.HandleFunc(
		"/follow",
		followHandler.FollowUser,
	)

	http.HandleFunc(
		"/following",
		followHandler.GetFollowing,
	)

	http.HandleFunc(
		"/tweets",
		tweetHandler.Tweets,
	)

	http.HandleFunc(
		"/timeline",
		timelineHandler.GetTimeline,
	)

	log.Println("🚀 server running :8080")

	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)

}
