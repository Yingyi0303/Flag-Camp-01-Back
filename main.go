package main

import (
	"fmt"
	"log"
	"net/http"

	"propertyManagement/backend"
	"propertyManagement/handler"
)

func main() {
    fmt.Println("Server is running on port 8080")

	backend.Init()
	defer backend.Close()

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
