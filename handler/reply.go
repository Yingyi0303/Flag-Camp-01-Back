package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"propertyManagement/model"
	"propertyManagement/service"

	"github.com/form3tech-oss/jwt-go"
)

func postReplyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one post reply request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	decoder := json.NewDecoder(r.Body)
	var reply model.Reply
	reply.Username = username.(string)
	if err := decoder.Decode(&reply); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	success, err := service.ValidateResidentialUser(username.(string))
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
	if reply.Content == "" {
		http.Error(w, "Invalid content", http.StatusBadRequest)
		return
	}

	// process request
	result, err := service.AddReply(&reply)
	if err != nil {
		http.Error(w, "Failed to add reply", http.StatusInternalServerError)
		fmt.Printf("Failed to add reply %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
	if result == nil {
		w.Write([]byte("{}"))
	} else {
		jsonResponse, _ := json.Marshal(result)
		w.Write(jsonResponse)
		fmt.Printf("Handler add reply: %d\n", reply.DiscussionId)
	}
}

func getMyRepliesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one get my replies request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	// validate request
	success, err := service.ValidateResidentialUser(username.(string))
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
	result, err := service.GetMyReplies(username.(string))
	if err != nil {
		http.Error(w, "Failed to get my replies", http.StatusInternalServerError)
		fmt.Printf("Failed to get my replies %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if len(result) == 0 {
		w.Write([]byte("[]"))
	} else {
		jsonResponse, _ := json.Marshal(result)
		w.Write(jsonResponse)
	}
	fmt.Printf("Handler get my replies\n")
}

func deleteReplyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one delete reply request")
	w.Header().Set("Content-Type", "application/json")

	// parse request
	token := r.Context().Value("user")
	claims := token.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	decoder := json.NewDecoder(r.Body)
	var reply model.Reply
	if err := decoder.Decode(&reply); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		fmt.Printf("Invalid input %v\n", err)
		return
	}

	// validate request
	success, err := service.ValidateResidentialUser(username.(string))
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
	err = service.RemoveReply(username.(string), reply.Id)
	if err != nil {
		http.Error(w, "Failed to remove reply", http.StatusInternalServerError)
		fmt.Printf("Failed to remove reply %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Printf("Handler delete reply : %d\n", reply.Id)
}
