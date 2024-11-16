package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tanmay958/app-docker/controller"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/api/submit",controller.SubmitJob ).Methods("POST")
	r.HandleFunc("/api/job/{jobID}", controller.GetJobDetails).Methods("GET") 
	// r.HandleFunc("/api/status", controller.GetJobStatus).Methods("GET")
	// Add more routes as needed...

	return r
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to the API Gateway!</h1>"))
}