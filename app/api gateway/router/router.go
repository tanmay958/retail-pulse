package router

import (
	"apigateway/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/",serveHome).Methods("GET")
    r.HandleFunc("/api/submit",controller.SubmitJob).Methods("POST")
    
    r.HandleFunc("/api/status", controller.GetJobStatus).Methods("GET")

    return r
}
func serveHome(w http.ResponseWriter , r *http.Request) {
	w.Write([]byte("<h1>hye Tanmay</h1>"))
}
