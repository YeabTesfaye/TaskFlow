package routes

import (
	"api/controllers"
	"api/middleware"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router) {
	// Auth routes
	r.HandleFunc("/api/users", controllers.SignUp).Methods("POST")
	r.HandleFunc("/api/login", controllers.LoginUser).Methods("POST")

	// User profile routes
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.GetMe)).
		Methods("GET")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.UpdateMe)).
		Methods("PATCH")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(controllers.DeleteMe)).
		Methods("DELETE")
	r.HandleFunc("/api/users/password", middleware.AuthMiddleware(
		controllers.ChangePassword)).Methods("POST")

	// Profile picture routes
	r.HandleFunc("/api/users/profile-picture", middleware.
		AuthMiddleware(controllers.UpdateProfilePicture)).Methods("POST")
	// User preferences routes
	r.HandleFunc("/api/users/preferences", middleware.
		AuthMiddleware(controllers.UpdatePreferences)).Methods("PUT")
	// Email verification routes
	r.HandleFunc("/api/users/send-verification", middleware.
		AuthMiddleware(controllers.SendVerificationEmail)).Methods("POST")
	r.HandleFunc("/api/users/verify-email", controllers.VerifyEmail).
		Methods("GET")
}
