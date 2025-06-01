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
	r.HandleFunc("/api/tasks/{id}/status", middleware.AuthMiddleware(
			controllers.UpdateTaskStatus)).Methods("PATCH")
	r.HandleFunc("/api/tasks/collaborators/add", middleware.AuthMiddleware(
		controllers.AddCollaborator)).Methods("POST")
	r.HandleFunc("/api/tasks/collaborators/remove", middleware.AuthMiddleware(
		controllers.RemoveCollaborator)).Methods("DELETE")


	// Comment routes
	r.HandleFunc("/api/tasks/{taskId}/comments", middleware.AuthMiddleware(
		controllers.AddComment)).Methods("POST")
	r.HandleFunc("/api/tasks/{taskId}/comments", middleware.AuthMiddleware(
		controllers.GetComments)).Methods("GET")
	r.HandleFunc("/api/tasks/{taskId}/comments/{commentId}", middleware.AuthMiddleware(
		controllers.UpdateComment)).Methods("PUT")
	r.HandleFunc("/api/tasks/{taskId}/comments/{commentId}", middleware.AuthMiddleware(
		controllers.DeleteComment)).Methods("DELETE")

	r.HandleFunc("/api/tags", middleware.AuthMiddleware(
		controllers.GetTags))


}
