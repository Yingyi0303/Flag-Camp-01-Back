package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"propertyManagement/model"
	"propertyManagement/service"

	"github.com/form3tech-oss/jwt-go"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
	w.Header().Set("Content-Type", "text/plain")

	// parse request
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	if user.Username == "" || regexp.MustCompile(`^[a-zA-Z0-9]$`).MatchString(user.Username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	if user.Role != "resident" && user.Role != "manager" && user.Role != "third_party" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// process request
	success, err := service.AddUser(&user)
	if err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		fmt.Printf("Failed to add user %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		fmt.Println("User already exists")
		return
	}
	
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Handler sign up user: %s\n", user.Username)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	w.Header().Set("Content-Type", "text/plain")

	// parse request
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	result, err := service.CheckUser(&user)
	if err != nil {
		http.Error(w, "Failed to check user", http.StatusInternalServerError)
		fmt.Printf("Failed to check user %v\n", err)
		return
	}

	if result == nil {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		fmt.Println("Wrong username or password")
		return
	}

	// process request
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":	result.Username,
		"role":	result.Role,
		"exp":	time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token %v\n", err)
		return
	}

	// construct response
	respone := model.Response {
		Username: result.Username,
		Role: result.Role,
		Token: tokenString,
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(respone)
	w.Write(jsonResponse)
	fmt.Printf("Handler sign in user: %s\n", user.Username)
}
