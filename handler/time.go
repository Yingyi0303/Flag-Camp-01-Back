package handler

import (
	"fmt"
	"net/http"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received time request")

	// process request
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte(formattedTime))
}
