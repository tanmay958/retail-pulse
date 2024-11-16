package controller

import (
	model "apigateway/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

	var stores   = make(map[string]*model.Store)
	var jobs     = make(map[string]*model.Job)
    var images = make(map[string]*model.Image)

func SubmitJob(w http.ResponseWriter, r *http.Request) {
   var requestData model.VisitsResponse
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	now := time.Now()
	jobID := now.Format("20060102150405") 
	for _, visit := range requestData.Visits {
		// Update Store objects
        fmt.Printf("%s\n",visit.StoreID )
		store, ok := stores[visit.StoreID]

		if !ok {
			store = &model.Store{
				StoreID:     visit.StoreID,
				VisitorCount: 0,
				Images:    make(map[string]model.Image),

			}
			stores[visit.StoreID] = store
		}
		store.VisitorCount++
        
        // stores[visit.StoreID] = 

		// Create Job objects
	
		job, ok := jobs[jobID]
		if !ok {
			job = &model.Job{
				JobID:   jobID,
				Images:  make(map[string]model.Image),
                Status:  "pending",
			}
			jobs[jobID] = job
		}

		for index, url := range visit.ImageURL {

			imageID := generateImageID() +  strconv.Itoa(index) +  visit.StoreID 
            image,ok :=  images[imageID]
            if !ok {
                image = &model.Image{
                    ImageID:imageID,
                    ImageURL:url ,
                    Status:"pending",

                }

                images[imageID] = image
            }
			job.Images[imageID] = *image
			store.Images[imageID] = *image
		}

		// Print job list
		
	}
    
        fmt.Println("current image list")
        for imageID, image := range images {
            fmt.Printf("  - Image ID: %s, URL: %s, Status: %s\n", imageID, image.ImageURL, image.Status)
        }
        fmt.Println("--------------")
        
	// Encode the response as JSON
    response := map[string]string{
		"jobid": jobID,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response) // Encode the response map
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
func generateImageID() string {
	return time.Now().Format("20060102150405999999")
}

func processJob(jobID uint) {
}


func GetJobStatus(w http.ResponseWriter, r *http.Request) {
    jobIDStr := r.URL.Query().Get("jobid")

	job, ok := jobs[jobIDStr] // Access the job from the jobs map
	if !ok {
		http.Error(w, `{"error":"Job not found"}`, http.StatusBadRequest)
		return
	}

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
	}

	// Prepare the response
	response := map[string]interface{}{
		"status": job.Status,
		"job_id": job.JobID,
	}

	// Include image details in the response
	imageDetails := make(map[string]map[string]interface{})
	for imageID, image := range job.Images {
		imageDetails[imageID] = map[string]interface{}{
			"status":    image.Status,
			"url":       image.ImageURL,
		}
	}
	response["images"] = imageDetails

	json.NewEncoder(w).Encode(response)
}
