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

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	http.HandleFunc("/users", userHandler.CreateUser)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
