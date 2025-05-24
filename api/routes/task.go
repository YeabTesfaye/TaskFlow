package routes

import (
	"api/controllers"
	"api/middleware"

	"github.com/gorilla/mux"
)

func RegisterTaskRoutes(r *mux.Router) {

	// Task statistics routes
	r.HandleFunc("/api/tasks/statistics", middleware.AuthMiddleware(
		controllers.GetTaskStatistics)).Methods("GET")
	r.HandleFunc("/api/tasks/statistics/by-category", middleware.AuthMiddleware(
		controllers.GetTaskStatisticsByCategory)).Methods("GET")

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

	// Comment routes
	r.HandleFunc("/api/tasks/{taskId}/comments", middleware.AuthMiddleware(
		controllers.AddComment)).Methods("POST")
	r.HandleFunc("/api/tasks/{taskId}/comments", middleware.AuthMiddleware(
		controllers.GetComments)).Methods("GET")
	r.HandleFunc("/api/tasks/{taskId}/comments/{commentId}", middleware.AuthMiddleware(
		controllers.UpdateComment)).Methods("PUT")
	r.HandleFunc("/api/tasks/{taskId}/comments/{commentId}", middleware.AuthMiddleware(
		controllers.DeleteComment)).Methods("DELETE")

	// Category routes
	r.HandleFunc("/api/categories", middleware.AuthMiddleware(
		controllers.CreateCategory)).Methods("POST")
	r.HandleFunc("/api/categories", middleware.AuthMiddleware(
		controllers.GetCategories)).Methods("GET")
	r.HandleFunc("/api/categories/{id}", middleware.AuthMiddleware(
		controllers.UpdateCategory)).Methods("PUT")
	r.HandleFunc("/api/categories/{id}", middleware.AuthMiddleware(
		controllers.DeleteCategory)).Methods("DELETE")
	r.HandleFunc("/api/tasks/category/{categoryId}", middleware.AuthMiddleware(
		controllers.GetTasksByCategory)).Methods("GET")

}
