package main

import (
	"api/configs"
	"api/jobs"
	"api/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Connect to MongoDB
	configs.ConnectDB()

	// Start the reminder job
	jobs.StartReminderJob()

	// Create router
	r := mux.NewRouter()

	// Register routes explicitly
	routes.RegisterUserRoutes(r)
	routes.RegisterTaskRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
