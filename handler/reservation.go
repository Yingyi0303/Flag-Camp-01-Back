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

func postReservationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post reservation request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var reservation model.Reservation
	reservation.Username = username
	if err := decoder.Decode(&reservation); err != nil {
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
	if reservation.FacilityName == "" {
		http.Error(w, "Invalid facility name", http.StatusBadRequest)
		return
	}
	if reservation.Remark == "" {
		http.Error(w, "Invalid remark", http.StatusBadRequest)
		return
	}
	reservationDate, err := time.Parse("2006-01-02", reservation.ReservationDate)
	if err != nil {
		http.Error(w, "Invalid reservation date", http.StatusBadRequest)
		return
	}
	currentDate, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if reservationDate.Before(currentDate) {
		http.Error(w, "Reservation for past date", http.StatusBadRequest)
		return
	}
	if reservation.StartHour < 0 {
		http.Error(w, "Invalid start hour", http.StatusBadRequest)
		return
	}
	if reservation.EndHour > 24 {
		http.Error(w, "Invalid end hour", http.StatusBadRequest)
		return
	}
	if reservation.EndHour - reservation.StartHour > 4 {
		http.Error(w, "Reservation longer than 4 hours", http.StatusBadRequest)
		return
	}
	if reservation.EndHour - reservation.StartHour < 1 {
		http.Error(w, "Reservation shorter than 1 hour", http.StatusBadRequest)
		return
	}

	// process request
	result, err := service.AddReservation(&reservation)
	if err != nil {
		http.Error(w, "Failed to add reservation", http.StatusInternalServerError)
		fmt.Printf("Failed to add reservation %v\n", err)
		return
	}

	if result == nil {
		http.Error(w, "Calendar schedule conflict", http.StatusInternalServerError)
		fmt.Println("Calendar schedule conflict")
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("Handler post reservation: %d\n", result.Id)
}

func getMyReservationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get my reservations request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	// validate request
	if role != "resident" && role != "manager" {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Println("User unauthorized")
		return
	}

	// process request
	result, err := service.GetMyReservations(username)
	if err != nil {
		http.Error(w, "Failed to get my reservations", http.StatusInternalServerError)
		fmt.Printf("Failed to get my reservations %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Handler get my reservations")
}

func deleteReservationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one delete reservation request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var reservation model.Reservation
	if err := decoder.Decode(&reservation); err != nil {
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

	// process request
	err := service.RemoveReservation(username, reservation.Id)
	if err != nil {
		http.Error(w, "Failed to remove reservation", http.StatusInternalServerError)
		fmt.Printf("Failed to remove reservation %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("Handler delete reservation: %d\n", reservation.Id)
}
