package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"propertyManagement/model"
	"propertyManagement/service"

	jwt "github.com/form3tech-oss/jwt-go"
)

func postDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post discussion request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var discussion model.Discussion
	discussion.Username = username
	if err := decoder.Decode(&discussion); err != nil {
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
	if discussion.Subject == "" {
		http.Error(w, "Invalid subject", http.StatusBadRequest)
		return
	}
	if discussion.Content == "" {
		http.Error(w, "Invalid content", http.StatusBadRequest)
		return
	}

	// process request
	result, err := service.AddDiscussion(&discussion)
	if err != nil {
		http.Error(w, "Failed to add discussion", http.StatusInternalServerError)
		fmt.Printf("Failed to add discussion %v\n", err)
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("Handler post discussion: %d\n", result.Id)
}

func getAllDiscussionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get all discussions request")
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
	result, err := service.GetAllDiscussions()
	if err != nil {
		http.Error(w, "Failed to get all discussions", http.StatusInternalServerError)
		fmt.Printf("Failed to get all discussions %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Handler get all discussions")
}

func getMyDiscussionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get my discussions request")
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
	result, err := service.GetMyDiscussions(username)
	if err != nil {
		http.Error(w, "Failed to get my discussions", http.StatusInternalServerError)
		fmt.Printf("Failed to get my discussions %v\n", err)
		return
	}

	jsonResponse, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Handler get my discussions")
}

func getDiscussionDetailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get discussion detail request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	// username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var discussion model.Discussion
	if err := decoder.Decode(&discussion); err != nil {
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
	result, err := service.GetDiscussionDetails(discussion.Id)
	if err != nil {
		http.Error(w, "Failed to get discussion details", http.StatusInternalServerError)
		fmt.Printf("Failed to get discussion details %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if result == nil {
		w.Write([]byte("{}"))
		return
	}
	jsonResponse, _ := json.Marshal(result)
	w.Write(jsonResponse)
	fmt.Printf("Handler get discussion details: %d\n", discussion.Id)
}

func deleteDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one delete discussion request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	decoder := json.NewDecoder(r.Body)
	var discussion model.Discussion
	if err := decoder.Decode(&discussion); err != nil {
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
	err := service.RemoveDiscussion(username, discussion.Id)
	if err != nil {
		http.Error(w, "Failed to remove discussion", http.StatusInternalServerError)
		fmt.Printf("Failed to remove discussion %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("Handler delete discussion: %d\n", discussion.Id)
}
