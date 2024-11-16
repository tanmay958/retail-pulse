package controller

import (
	model "apigateway/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)
var (
	jobCounter int
	jobMutex   sync.Mutex
)
var jobs []model.Job


func SubmitJob(w http.ResponseWriter, r *http.Request) {
   var requestData model.VisitsResponse
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, visit := range requestData.Visits {
		jobMutex.Lock()
		jobCounter++
		job := model.Job{
			JobID:    jobCounter,
			StoreID:  visit.StoreID,
			ImageURL: visit.ImageURL,
		}
		jobs = append(jobs, job)
		jobMutex.Unlock()
	}
    fmt.Println("Current job list:")
		for _, j := range jobs {
			fmt.Printf("  - Job ID: %d, Store ID: %s, Image URLs: %v\n", j.JobID, j.StoreID, j.ImageURL)
		}

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(requestData)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func processJob(jobID uint) {
    // var job models.Job
    // database.DB.Preload("Visits.ImageURLs").First(&job, jobID)

    // for _, visit := range job.Visits {
    //     for _, image := range visit.ImageURLs {
    //         time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond) // Simulate processing
    //         perimeter := calculatePerimeter(image.URL)
    //         database.DB.Model(&image).Update("Perimeter", perimeter)
    //     }
    // }

    // database.DB.Model(&job).Update("Status", "completed")
}

func calculatePerimeter(url string) float64 {
    // Mock perimeter calculation (replace with real logic)
    return 2 * (1920 + 1080) // Assuming a fixed resolution
}

func GetJobStatus(w http.ResponseWriter, r *http.Request) {
    // jobIDStr := r.URL.Query().Get("jobid")
    // jobID, err := strconv.Atoi(jobIDStr)
    // if err != nil {
    //     http.Error(w, `{"error":"Invalid job ID"}`, http.StatusBadRequest)
    //     return
    // }

    // var job models.Job
    // if err := database.DB.First(&job, jobID).Error; err != nil {
    //     http.Error(w, `{"error":"Job not found"}`, http.StatusBadRequest)
    //     return
    // }

    // if job.Status == "failed" {
    //     var visits []models.Visit
    //     database.DB.Where("job_id = ?", job.ID).Find(&visits)

    //     failedStores := []map[string]string{}
    //     for _, visit := range visits {
    //         failedStores = append(failedStores, map[string]string{"store_id": visit.StoreID})
    //     }

    //     json.NewEncoder(w).Encode(map[string]interface{}{
    //         "status": "failed",
    //         "job_id": job.ID,
    //         "error":  failedStores,
    //     })
    //     return
    // }

    // json.NewEncoder(w).Encode(map[string]interface{}{
    //     "status": job.Status,
    //     "job_id": job.ID,
    // })
}
