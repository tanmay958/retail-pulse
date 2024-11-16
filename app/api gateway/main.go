package main

import (
	"apigateway/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
    // Initialize database
	fmt.Println("server starting ")
    // database.InitDatabase()

    // Setup router
    r := router.SetupRouter()

    // Start server
    log.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
