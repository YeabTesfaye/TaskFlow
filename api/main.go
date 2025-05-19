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
	r.HandleFunc("/api/users", controllers.SignUp).
		Methods("POST")
	r.HandleFunc("/api/login", controllers.LoginUser).
		Methods("POST")

	// Protected routes
	r.HandleFunc("/api/tasks", middleware.
		AuthMiddleware(controllers.CreateTask)).Methods("POST")
	r.HandleFunc("/api/tasks", middleware.
		AuthMiddleware(controllers.GetUserTasks)).Methods("GET")
	r.HandleFunc("/api/tasks/{id}", middleware.
		AuthMiddleware(controllers.UpdateTask)).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", middleware.
		AuthMiddleware(controllers.DeleteTask)).Methods("DELETE")
	r.HandleFunc("/api/tasks/{id}", middleware.AuthMiddleware(controllers.GetTask)).
		Methods("GET")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.DeleteMe)).
		Methods("DELETE")
		// ... existing code ...
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.UpdateMe)).
		Methods("PATCH")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.GetMe)).
		Methods("GET")
	r.HandleFunc("/api/users/password", middleware.
		AuthMiddleware(controllers.ChangePassword)).Methods("POST")
	r.HandleFunc("/api/tasks/stats", middleware.
	AuthMiddleware(controllers.GetTaskStats)).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
