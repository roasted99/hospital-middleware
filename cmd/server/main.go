package main

import (
  "log"
  "net/http"
  "os"
  "fmt"

  "github.com/gorilla/mux"
  "github.com/roasted99/hospital-middleware/internal/config"
  "github.com/roasted99/hospital-middleware/internal/db"
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

  // Define routes
  // router.HandleFunc("/api/v1/staff/login", loginHandler).Methods("POST")
  // router.HandleFunc("/api/v1/staff/register", registerHandler).Methods("POST")

  // Start server
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  
  fmt.Printf("Server is running on port %s\n", port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}