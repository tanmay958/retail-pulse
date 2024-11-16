package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tanmay958/app-docker/config"
	model "github.com/tanmay958/app-docker/models"

	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error

	// Connect to MySQL
	db, err = config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Connected to the database successfully!")

	// Migrate models

	err = db.AutoMigrate(&model.Store{}, &model.Image{}, &model.Job{}, &model.Visit{})
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	fmt.Println("Migration completed successfully!")

	// Start HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello, the database is connected and ready!")
	})
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
