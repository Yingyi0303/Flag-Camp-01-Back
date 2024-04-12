package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"propertyManagement/model"
	"propertyManagement/service"

	"github.com/form3tech-oss/jwt-go"
)

func postMaintenanceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post maintenance request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)

	decoder := json.NewDecoder(r.Body)
	var maintenance model.Maintenance
	maintenance.Username = username
	if err := decoder.Decode(&maintenance); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	success, err := service.ValidateResidentialUser(username)
	if err != nil {
		http.Error(w, "Failed to validate role", http.StatusInternalServerError)
		fmt.Printf("Failed to validate role %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Printf("User unauthorized %v\n", err)
		return
	}
	if maintenance.Subject == "" {
		http.Error(w, "Invalid title", http.StatusBadRequest)
		return
	}
	if maintenance.Content == "" {
		http.Error(w, "Invalid content", http.StatusBadRequest)
		return
	}

	// process request
	result, err := service.AddMaintenance(&maintenance)
	if err != nil {
		http.Error(w, "Failed to add maintenance", http.StatusInternalServerError)
		fmt.Printf("Failed to add maintenance %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("Handler post maintenance: %d\n", result.Id)
}

func getAllMaintenancesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get all maintenances request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)

	decoder := json.NewDecoder(r.Body)
	var maintenance model.Maintenance
	maintenance.Username = username
	if err := decoder.Decode(&maintenance); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	success, err := service.ValidateThirdPartyUser(username)
	if err != nil {
		http.Error(w, "Failed to validate role", http.StatusInternalServerError)
		fmt.Printf("Failed to validate role %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Printf("User unauthorized %v\n", err)
		return
	}

	// process request
	result, err := service.GetAllMaintenances(maintenance.Completed)
	if err != nil {
		http.Error(w, "Failed to get all maintenances", http.StatusInternalServerError)
		fmt.Printf("Failed to get all maintenances %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("Handler get all maintenances\n")
}

func getMyMaintenancesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get my maintenances request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)

	decoder := json.NewDecoder(r.Body)
	var maintenance model.Maintenance
	maintenance.Username = username
	if err := decoder.Decode(&maintenance); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	success, err := service.ValidateResidentialUser(username)
	if err != nil {
		http.Error(w, "Failed to validate role", http.StatusInternalServerError)
		fmt.Printf("Failed to validate role %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Printf("User unauthorized %v\n", err)
		return
	}

	// process request
	result, err := service.GetMyMaintenances(maintenance.Username, maintenance.Completed)
	if err != nil {
		http.Error(w, "Failed to get my maintenances", http.StatusInternalServerError)
		fmt.Printf("Failed to get my maintenances %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("Handler get my maintenances\n")
}

func putMaintenanceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one put maintenances request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)

	decoder := json.NewDecoder(r.Body)
	var maintenance model.Maintenance
	maintenance.Username = username
	if err := decoder.Decode(&maintenance); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	success, err := service.ValidateThirdPartyUser(username)
	if err != nil {
		http.Error(w, "Failed to validate role", http.StatusInternalServerError)
		fmt.Printf("Failed to validate role %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		fmt.Printf("User unauthorized %v\n", err)
		return
	}

	// process request
	result, err := service.SetMaintenance(&maintenance)
	if err != nil {
		http.Error(w, "Failed to set maintenance", http.StatusInternalServerError)
		fmt.Printf("Failed to set maintenance %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if result == nil {
		w.Write([]byte("{}"))
		return
	}
	jsonResponse, _ := json.Marshal(result)
	w.Write(jsonResponse)
	fmt.Printf("Handler put maintenance: %d\n", maintenance.Id)	
}
