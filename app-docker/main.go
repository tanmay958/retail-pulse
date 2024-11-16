package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tanmay958/app-docker/config"
	"github.com/tanmay958/app-docker/controller"
	model "github.com/tanmay958/app-docker/models"
	"github.com/tanmay958/app-docker/router"
	"github.com/tanmay958/app-docker/utils"

	"gorm.io/gorm"
)

var db *gorm.DB
func init()  {
	var err error

	// Connect to MySQL
	db, err = config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Connected to the database successfully!")

	// Migrate models

	db.AutoMigrate(&model.Visit{}, &model.Job{}, &model.Store{}, &model.Image{})
	if err != nil {
		log.Fatal("migration failed:", err)
	}
	controller.SetDB(db) 
	fmt.Println("Migration completed successfully!")

	// Start HTTP server

}

func main() {
	csvFilePath := filepath.Join("/app", "assets", "static", "StoreMasterAssignment.csv")
	fmt.Println("started")
	utils.DumpStoresFromCSV(csvFilePath, db)
	fmt.Println("server starting ")
    // database.InitDatabase()

    // Setup router
    r := router.SetupRouter()

    // Start server
    log.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))


}

