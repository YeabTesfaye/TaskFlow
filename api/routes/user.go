package routes

import (
	"api/controllers"
	"api/middleware"

	"github.com/gorilla/mux"
)
func RegisterUserRoutes(r *mux.Router, userController *controllers.UserController, profileController *controllers.ProfileController) {
	// Auth routes
	r.HandleFunc("/api/signup", userController.SignUp).Methods("POST")
	r.HandleFunc("/api/login", userController.Login).Methods("POST")

	// User profile routes
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(userController.GetMe)).
		Methods("GET")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(userController.UpdateMe)).
		Methods("PATCH")
	r.HandleFunc("/api/users/me", middleware.AuthMiddleware(userController.DeleteMe)).
		Methods("DELETE")
	r.HandleFunc("/api/users/password", middleware.AuthMiddleware(
		userController.ChangePassword)).Methods("POST")

	// Profile picture routes
	r.HandleFunc("/api/users/profile-picture", middleware.
		AuthMiddleware(profileController.UpdateProfilePicture)).Methods("POST")

	// User preferences routes
	r.HandleFunc("/api/users/preferences", middleware.
		AuthMiddleware(profileController.UpdatePreferences)).Methods("PUT")

	// Email verification routes
	r.HandleFunc("/api/users/send-verification", middleware.
		AuthMiddleware(userController.SendVerificationEmail)).Methods("POST")
	r.HandleFunc("/api/users/verify-email", userController.VerifyEmail).
		Methods("GET")
}