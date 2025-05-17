package main

import (
	"api/configs"
	"api/controllers"
	"api/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// connect to MongoDB
	configs.ConnectDB()

	// Create router
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login", controllers.LoginUser).Methods("POST")

	// Protected routes
	r.HandleFunc("/tasks", middleware.
		AuthMiddleware(controllers.CreateTask)).Methods("POST")
	r.HandleFunc("/tasks", middleware.
	AuthMiddleware(controllers.GetUserTasks)).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
