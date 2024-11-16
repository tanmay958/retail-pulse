package main

import (
	"apigateway/router"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("server is started at port 8080")
	r := router.InitializeRouter()
	http.ListenAndServe(":8080", r)
	fmt.Println("Listening")
}