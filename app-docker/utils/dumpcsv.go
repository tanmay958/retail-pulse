package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	model "github.com/tanmay958/app-docker/models"
	"gorm.io/gorm"
)

func DumpStoresFromCSV(csvFilePath string, db *gorm.DB) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal("failed to open CSV file:", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("failed to read CSV:", err)
	}

	// Prepare a slice to hold Store records
	var stores []model.Store

	// Iterate through each record and map it to Store struct
	for _, record := range records {
		if len(record) < 3 {
			// Skip invalid rows (e.g., empty lines or incomplete records)
			continue
		}

		store := model.Store{
			StoreID:   record[0], // Store ID from CSV
			StoreName: record[1], // Store Name from CSV
			AreaCode:  record[2], // Area Code from CSV
		}

		// Append to stores slice
		stores = append(stores, store)
	}

	// Insert stores into the database in bulk
	if err := db.Create(&stores).Error; err != nil {
		log.Fatal("failed to insert stores into the database:", err)
	}
	fmt.Println("Stores have been successfully added to the database!")
}