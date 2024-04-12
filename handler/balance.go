package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"propertyManagement/service"

	"github.com/form3tech-oss/jwt-go"
)

func getMyBalanceHandlder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get my balance request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)

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
	result, err := service.GetMyBalance(username)
	if err != nil {
		http.Error(w, "Failed to get my balance", http.StatusInternalServerError)
		fmt.Printf("Failed to get my balance %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("Handler get my balance\n")
}
