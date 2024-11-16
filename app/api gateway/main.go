package main

import (
	"apigateway/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
  
	fmt.Println("server starting ")
    r := router.SetupRouter()
    log.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
