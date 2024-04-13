package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"propertyManagement/model"
	"propertyManagement/service"

	"github.com/form3tech-oss/jwt-go"
)

func postPaymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post payment request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var payment model.Payment
	payment.Username = username
	if err := decoder.Decode(&payment); err != nil {
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
	if payment.Item == "" {
		http.Error(w, "Invalid item", http.StatusBadRequest)
		return
	}
	if payment.Amount < 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// process request
	result, err := service.AddPayment(&payment)
	if err != nil {
		http.Error(w, "Failed to add payment", http.StatusInternalServerError)
		fmt.Printf("Failed to add payment %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if result == nil {
		w.Write([]byte("{}"))
		return
	}
	jsonResponse, _ := json.Marshal(result)
	w.Write(jsonResponse)
	fmt.Printf("Handler post payment: %d\n", result.Id)
}

func getMyPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get my payments request")
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
	result, err := service.GetMyPayments(username)
	if err != nil {
		http.Error(w, "Failed to get my payments", http.StatusInternalServerError)
		fmt.Printf("Failed to get my payments %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Handler get my payments")
}
