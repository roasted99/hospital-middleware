package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/roasted99/hospital-middleware/internal/config"
)

// InitDB initializes database connection
func InitDB() (*sql.DB, error) {
	dbConfig := config.GetDBConfig()
	
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)
	
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	
	if err = db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}