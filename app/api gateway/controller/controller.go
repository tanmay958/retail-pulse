package controller

import (
	model "apigateway/models"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
   
    imageMux sync.Mutex // Mutex for images map
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
                    
                    ImageURL:url ,
                    Status:"pending",

                }

                images[imageID] = image
            }
			job.Images[imageID] = *image
			store.Images[imageID] = *image
		}
      

		
		
	}
    go processJob(jobID)
    
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

func processJob(jobID string) {
    fmt.Println(jobID)
	job, ok := jobs[jobID]
	if !ok {
		fmt.Printf("Job %s not found\n", jobID)
		return
	}

	var wg sync.WaitGroup // WaitGroup for parallel processing
	for imageID, image := range job.Images {
		wg.Add(1) // Increment WaitGroup counter
  
		go func(imageID string, image *model.Image) {
			defer wg.Done() 

			
			fmt.Printf("Downloading image %s from URL %s\n", imageID, image.ImageURL)
			
			
			height := 1080 
			width := 1920  
			perimeter := 2 * (height + width)
           
            time.Sleep(time.Duration(rand.Float64()*2+5) * time.Second) 
			// time.Sleep(time.Duration(rand.Float64()*0.3+0.1) * time.Second)

			
			imageMux.Lock()
			image.Status = "completed"
			image.Perimeter = perimeter
			images[imageID] = image
            job.Images[imageID] =  *image
			imageMux.Unlock()

		}(imageID, &image)
	}
    
	wg.Wait() 
  
}

func GetJobStatus(w http.ResponseWriter, r *http.Request) {
    jobIDStr := r.URL.Query().Get("jobid")

	job, ok := jobs[jobIDStr] // Access the job from the jobs map
  
	if !ok {
		http.Error(w, `{"error":"Job not found"}`, http.StatusBadRequest)
		return
	}
    for imageID, image := range job.Images {
        fmt.Printf("  - Image ID: %s, URL: %s, Status: %s\n", imageID, image.ImageURL, image.Status)
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
