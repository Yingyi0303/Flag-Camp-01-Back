package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"propertyManagement/model"
	"propertyManagement/service"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

func getAllFacilitiesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get all facilities request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	// username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	// validate request
	if role != "resident" && role != "manager" {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Println("User unauthorized")
		return
	}

	// process request
	result, err := service.GetAllFacilities()
	if err != nil {
		http.Error(w, "Failed to get all facilities", http.StatusInternalServerError)
		fmt.Printf("Failed to get all facilities %v\n", err)
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Handler get all facilities")
}

func getFacilityReservationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get facility reservations request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	// username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var request model.ReservationRequest
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	if role != "resident" && role != "manager" {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Println("User unauthorized")
		return
	}
	if request.FacilityName == "" {
		http.Error(w, "Invalid facility name", http.StatusBadRequest)
		return
	}
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}
	currentDate, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if startDate.Before(currentDate) {
		http.Error(w, "Start date in the past", http.StatusBadRequest)
		return
	}
	if endDate.Before(currentDate) {
		http.Error(w, "End date in the past", http.StatusBadRequest)
		return
	}
	if endDate.Before(startDate) {
		http.Error(w, "End date before start date", http.StatusBadRequest)
		return
	}

	// process request
	result, err := service.GetFacilityReservations(&request)
	if err != nil {
		http.Error(w, "Failed to get reservations", http.StatusInternalServerError)
		fmt.Printf("Failed to get reservations %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("Handler get reservations for facility: %s\n", request.FacilityName)
}
