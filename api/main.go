package main

import (
	"api/configs"
	"api/controllers"
	"api/jobs"
	"api/middleware"
	"api/repositories"
	"api/routes"
	"api/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	userRepo := repositories.NewUserRepository(configs.GetCollection(configs.DB, "users"))
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	profileService := services.NewProfileService(userRepo) // Add this
	profileController := controllers.NewProfileController(profileService)
	// Connect to MongoDB
	configs.ConnectDB()

	// Start the reminder job
	jobs.StartReminderJob()

	// Create router
	r := mux.NewRouter()

	// Add rate limiting middleware
	r.Use(middleware.RateLimit)

	// Add input sanitize middleware
	r.Use(middleware.SanitizeInput)

	// Register routes explicitly
	routes.RegisterUserRoutes(r,userController, profileController)
	routes.RegisterTaskRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
