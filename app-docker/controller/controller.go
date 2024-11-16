package controller

import (
	"encoding/json"

	"io"
	"net/http"
	"time"

	model "github.com/tanmay958/app-docker/models"

	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the database instance for the controller package
func SetDB(database *gorm.DB) {
	db = database
}

func SubmitJob(w http.ResponseWriter, r *http.Request) {
	// Efficiently parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the JSON data into a generic map
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Extract the count and visits from the map
	count, ok := data["count"].(float64)
	if !ok {
		http.Error(w, `{"error":"Invalid 'count' field"}`, http.StatusBadRequest)
		return
	}

	visits, ok := data["visits"].([]interface{})
	if !ok {
		http.Error(w, `{"error":"Invalid 'visits' field"}`, http.StatusBadRequest)
		return
	}

	// Validate the count and visits
	if int(count) == 0 || len(visits) == 0 || int(count) != len(visits) {
		http.Error(w, `{"error":"Missing fields or count mismatch"}`, http.StatusBadRequest)
		return
	}

	// Database interaction starts here
	now := time.Now()
	jobID := now.Format("20060102150405")

	// Create a new job in the database
	job := model.Job{
		JobID:  jobID,
		Status: "pending",
	}
	result := db.Create(&job)
	if result.Error != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	for _, v := range visits {
		visit, ok := v.(map[string]interface{})
		if !ok {
			// Handle invalid visit data
			continue
		}

		storeID, ok := visit["store_id"].(string)
		if !ok {
			// Handle missing or invalid store_id
			continue
		}

		imageURLs, ok := visit["image_url"].([]interface{})
		if !ok {
			// Handle missing or invalid image_url
			continue
		}

		// Convert imageURLs to []string
		var imageURLStrings []string
		for _, url := range imageURLs {
			if str, ok := url.(string); ok {
				imageURLStrings = append(imageURLStrings, str)
			}
		}

		// Find or create the store in the database
		var store model.Store
		result := db.FirstOrCreate(&store, model.Store{StoreID: storeID})
		if result.Error != nil {
			http.Error(w, "Failed to find or create store", http.StatusInternalServerError)
			return
		}

		store.VisitorCount++
		result = db.Save(&store)
		if result.Error != nil {
			http.Error(w, "Failed to update store", http.StatusInternalServerError)
			return
		}

		for _, url := range imageURLStrings {
			// imageID := generateImageID() + strconv.Itoa(index) + storeID

			// Create a new image in the database
			image := model.Image{
				StoreID:  storeID,
				JobID:    jobID,
				ImageURL: url,
				Status:   "pending",
			}
			result := db.Create(&image)
			if result.Error != nil {
				http.Error(w, "Failed to create image", http.StatusInternalServerError)
				return
			}
		}
	}

	// go processJob(jobID)

	// Encode the response as JSON
	response := map[string]string{
		"jobid": jobID,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func generateImageID() string {
	return time.Now().Format("20060102150405999999")
}

// Placeholder for your actual job processing logic
// func processJob(jobID string) {
// 	// Implement your job processing logic here
// 	fmt.Println("Processing job:", jobID)
// 	// ... update image statuses, perform calculations, etc. ...
// }