package main

import (
	"fmt"
	"net/http"

	"twitter-system-design/internal/database"
	"twitter-system-design/internal/handlers"
	"twitter-system-design/internal/repositories"
	"twitter-system-design/internal/services"
)

func main() {
	db := database.NewPostgres()

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	followRepo := repositories.NewFollowRepository(db)
	tweetRepo := repositories.NewTweetRepositrory(db)

	// Services
	userService := services.NewUserService(userRepo)
	followService := services.NewFollowService(followRepo)
	tweetService := services.NewTweetService(tweetRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	followHandler := handlers.NewFollowHandler(followService)
	tweetHandler := handlers.NewTweetHandler(tweetService)

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)

	http.HandleFunc("/follow", followHandler.FollowUser)
	http.HandleFunc("/following", followHandler.GetFollowing)

	http.HandleFunc("/tweets", tweetHandler.CreateTweet)
	http.HandleFunc("/user-tweets", tweetHandler.GetTweetsByUserID)

	fmt.Println("🚀 Server running on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
