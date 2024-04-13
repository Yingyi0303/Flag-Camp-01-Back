package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"propertyManagement/model"
	"propertyManagement/service"

	"github.com/form3tech-oss/jwt-go"
)

func postBillHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post bill request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	// username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var bill model.Bill
	if err := decoder.Decode(&bill); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	if role != "third_party" {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Println("User unauthorized")
		return
	}
	if bill.Item == "" {
		http.Error(w, "Invalid item", http.StatusBadRequest)
		return
	}
	if bill.Amount < 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// process request
	result, err := service.AddBill(&bill)
	if err != nil {
		http.Error(w, "Failed to add bill", http.StatusInternalServerError)
		fmt.Printf("Failed to add bill %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if result == nil {
		w.Write([]byte("{}"))
		return
	}
	jsonResponse, _ := json.Marshal(result)
	w.Write(jsonResponse)
	fmt.Printf("Handler post bill: %d\n", result.Id)
}

func getMyBillsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get my bills request")
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
	result, err := service.GetMyBills(username)
	if err != nil {
		http.Error(w, "Failed to get my bills", http.StatusInternalServerError)
		fmt.Printf("Failed to get my bills %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Handler get my bills")
}
