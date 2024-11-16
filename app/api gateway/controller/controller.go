package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// RequestPayload represents the structure for job submission.
type RequestPayload struct {
	Count  int `json:"count"`
	Visits []struct {
		StoreID   string   `json:"store_id"`
		ImageURL  []string `json:"image_url"`
		VisitTime string   `json:"visit_time"`
	} `json:"visits"`
}

// JobStatus represents the response for job status.
type JobStatus struct {
	Status string      `json:"status"`
	JobID  int         `json:"job_id"`
	Error  interface{} `json:"error,omitempty"`
}

// Job represents an in-memory store for jobs.
type Job struct {
	ID      int
	Payload RequestPayload
	Status  string
	Error   []struct {
		StoreID string `json:"store_id"`
		Error   string `json:"error"`
	}
}

var jobIDCounter = 0
var jobs = make(map[int]Job)

// SubmitJobHandler handles the job submission.
func SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Validate request payload
	if payload.Count != len(payload.Visits) {
		http.Error(w, `{"error": "count does not match the number of visits"}`, http.StatusBadRequest)
		return
	}

	// Create job
	jobIDCounter++
	job := Job{
		ID:      jobIDCounter,
		Payload: payload,
		Status:  "ongoing",
	}
	jobs[jobIDCounter] = job

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"job_id": jobIDCounter})
}

// GetJobStatusHandler handles the job status query.
func GetJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	jobIDParam := r.URL.Query().Get("jobid")
	if jobIDParam == "" {
		http.Error(w, `{"error": "jobid is required"}`, http.StatusBadRequest)
		return
	}

	jobID, err := strconv.Atoi(jobIDParam)
	if err != nil || jobID <= 0 {
		http.Error(w, `{"error": "invalid jobid"}`, http.StatusBadRequest)
		return
	}

	job, exists := jobs[jobID]
	if !exists {
		http.Error(w, `{message :  No such jobId found}`, http.StatusBadRequest) // JobID not found
		return
	}

	
	statusResponse := JobStatus{
		Status: job.Status,
		JobID:  job.ID,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statusResponse)
}
