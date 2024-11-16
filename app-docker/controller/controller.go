package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	model "github.com/tanmay958/app-docker/models"
	"github.com/tanmay958/app-docker/utils"
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

	go processJob(jobID)

	// Encode the response as JSON
	response := map[string]string{
		"jobid": jobID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// func generateImageID() string {
// 	return time.Now().Format("20060102150405999999")
// }
func GetJobDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]

	// Fetch the job details from the database
	var job model.Job
	result := db.Preload("Images").First(&job, "job_id = ?", jobID) // Include images
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Job not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch job details", http.StatusInternalServerError)
		}
		return
	}

	// Encode the job details as JSON
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(job)
	if err != nil {
		http.Error(w, "Failed to encode job details", http.StatusInternalServerError)
		return
	}
}
// Placeholder for your actual job processing logic
// func processJob(jobID string) {
// 	// Implement your job processing logic here
// 	fmt.Println("Processing job:", jobID)
// 	// ... update image statuses, perform calculations, etc. ...
// }
func processJob(jobID string) {
	fmt.Println("Processing job:", jobID)

	// Fetch the job with its associated images from the database
	var job model.Job
	result := db.Preload("Images").First(&job, "job_id = ?", jobID)
	if result.Error != nil {
		fmt.Printf("Job %s not found or error fetching: %v\n", jobID, result.Error)
		return
	}

	var wg sync.WaitGroup
	for _, image := range job.Images {
		wg.Add(1)
		go func(image model.Image) {
			defer wg.Done()

			time.Sleep(time.Duration(rand.Float64()*2+10) * time.Second)

			perimeter, err := utils.CalculatePerimeter(image.ImageURL)
			if err != nil {
				// Update image status to "failed" and store the error message
				image.Status = "failed"
				image.ErrMessage = err.Error()
			} else {
				// Update image status to "completed" and store the perimeter
				image.Status = "completed"
				image.Perimeter = perimeter
			}

			// Save the updated image in the database
			result := db.Save(&image)
			if result.Error != nil {
				fmt.Printf("Error updating image %d: %v\n", image.ID, result.Error)
			}
		}(image)
	}

	wg.Wait()
}

func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.URL.Query().Get("jobid")

	// Fetch the job with its associated images from the database
	var job model.Job
	result := db.Preload("Images").First(&job, "job_id = ?", jobIDStr)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, `{"error":"Job not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch job details", http.StatusInternalServerError)
		}
		return
	}

	// Print image details (for debugging)
	for _, image := range job.Images {
		fmt.Printf(" - Image ID: %d, URL: %s, Status: %s\n", image.ID, image.ImageURL, image.Status)
	}
	fmt.Println("--------------")

	// Calculate the overall job status based on image statuses
	allCompleted := true
	for _, image := range job.Images {
		if image.Status != "completed" {
			allCompleted = false
			break
		}
	}

	// Update job status if all images are completed
	if allCompleted {
		job.Status = "completed"
		db.Save(&job) // Update the job status in the database
	}

	failedImages := []map[string]string{}
	for _, image := range job.Images {
		if image.Status == "failed" {
			failedImages = append(failedImages, map[string]string{
				"image_id": fmt.Sprintf("%d", image.ID), // Use image ID from database
				"error":    image.ErrMessage,
			})
		}
	}

	if len(failedImages) > 0 {
		response := map[string]interface{}{
			"status": "failed",
			"job_id": job.JobID,
			"error":  failedImages,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"status": job.Status,
		"job_id": job.JobID,
	}

	json.NewEncoder(w).Encode(response)
}