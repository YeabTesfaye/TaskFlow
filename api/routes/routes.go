package routes

import (
	"api/controllers"
	"api/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes(r *mux.Router) {
	// Add error handler middleware
	r.Use(middleware.ErrorHandler)

	// Setup routes by category
	setupPublicRoutes(r)
	setupUserRoutes(r)
	setupTaskRoutes(r)
}

// setupPublicRoutes configures routes that don't require authentication
func setupPublicRoutes(r *mux.Router) {
	// Auth routes
	r.HandleFunc("/api/users", controllers.SignUp).Methods("POST")
	r.HandleFunc("/api/login", controllers.LoginUser).Methods("POST")
}

// setupUserRoutes configures user-related protected routes
func setupUserRoutes(r *mux.Router) {
	// User management routes
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.GetMe)).
		Methods("GET")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.UpdateMe)).
		Methods("PATCH")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.DeleteMe)).
		Methods("DELETE")
	r.HandleFunc("/api/users/password", middleware.AuthMiddleware(controllers.ChangePassword)).
		Methods("POST")
}

// setupTaskRoutes configures task-related protected routes
func setupTaskRoutes(r *mux.Router) {
	// Task utility routes
	r.HandleFunc("/api/tasks/stats", middleware.AuthMiddleware(controllers.GetTaskStats)).
		Methods("GET")

	// Task management routes
	r.HandleFunc("/api/tasks", middleware.AuthMiddleware(controllers.CreateTask)).
		Methods("POST")
	r.HandleFunc("/api/tasks", middleware.AuthMiddleware(controllers.GetUserTasks)).
		Methods("GET")
	r.HandleFunc("/api/tasks/{id}", middleware.AuthMiddleware(controllers.GetTask)).
		Methods("GET")
	r.HandleFunc("/api/tasks/{id}", middleware.AuthMiddleware(controllers.UpdateTask)).
		Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", middleware.AuthMiddleware(controllers.DeleteTask)).
		Methods("DELETE")
}
