package main

import (
  "log"
  "net/http"
  "os"
  "fmt"

  "github.com/gorilla/mux"
  "github.com/roasted99/hospital-middleware/internal/config"
  "github.com/roasted99/hospital-middleware/internal/db"
  "github.com/roasted99/hospital-middleware/internal/api/handlers"
  "github.com/roasted99/hospital-middleware/internal/api/middleware"
)

func main() {
  // Load environment variables
  if err := config.Load(); err != nil {
    log.Fatalf("Error loading .env file: %v", err)
  }

  // Initialize database connection
  db, err := db.InitDB()
  if err != nil {
    log.Fatalf("Error initializing database: %v", err)
  }
  defer db.Close()

  // Initialize router
  router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/staff/create", handlers.CreateStaff(db)).Methods("POST")
	router.HandleFunc("/staff/login", handlers.LoginStaff(db)).Methods("POST")
	
	// Protected routes
	patientRouter := router.PathPrefix("/patient").Subrouter()
	patientRouter.Use(middleware.Authenticate)
	patientRouter.HandleFunc("/search", handlers.SearchPatient(db)).Methods("GET")  

  // Start server
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  
  fmt.Printf("Server is running on port %s\n", port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}