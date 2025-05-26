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
	"github.com/rs/cors"
)

func main() {
	// Initialize repositories, services, controllers
	userRepo := repositories.NewUserRepository(configs.GetCollection(configs.DB, "users"))
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	profileService := services.NewProfileService(userRepo)
	profileController := controllers.NewProfileController(profileService)

	// Connect to MongoDB
	configs.ConnectDB()

	// Start background jobs
	jobs.StartReminderJob()

	// Create router
	r := mux.NewRouter()

	// Add middlewares
	r.Use(middleware.RateLimit)
	r.Use(middleware.SanitizeInput)

	// Register your routes
	routes.RegisterUserRoutes(r, userController, profileController)
	routes.RegisterTaskRoutes(r)

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Wrap router with CORS
	handler := corsHandler.Handler(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
