package router

import (
	"apigateway/controller"

	"github.com/gorilla/mux"
)

func InitializeRouter() *mux.Router {
	r := mux.NewRouter()

	// Route to submit a job
	r.HandleFunc("/api/submit/", controller.SubmitJobHandler).Methods("POST")

	// Route to get job info
	r.HandleFunc("/api/status", controller.GetJobStatusHandler).Methods("GET")

	return r
}