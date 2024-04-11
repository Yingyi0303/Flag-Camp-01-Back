package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func AddUser(user *model.User) (bool, error) {
	exists, err := backend.PGBackend.UserExists(user.Username)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if exists {
		return false, nil
	}

	err = backend.PGBackend.InsertUser(user.Username, user.Password, user.Role)
	if err != nil {
		log.Println(err)
		return false, err
	}
	fmt.Printf("Service added user: %s\n", user.Username)
	return true, nil
}

func CheckUser(user *model.User) (bool, error) {
	exists, err := backend.PGBackend.ValidateUser(user.Username, user.Password)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if !exists {
		return false, nil
	}

	fmt.Printf("Service chekced user: %s\n", user.Username)
	return true, nil
}

func ValidateThirdPartyUser(username string) (bool, error) {
	role, err := backend.PGBackend.ValidateRole(username)
	if err != nil {
		log.Println(err)
		return false, err
	}

	fmt.Printf("Service chekced third party: %s\n", username)
	if role == "third_party" {
		return true, nil
	} else {
		return false, nil
	}
}

func ValidateResidentialUser(username string) (bool, error) {
	role, err := backend.PGBackend.ValidateRole(username)
	if err != nil {
		log.Println(err)
		return false, err
	}

	fmt.Printf("Service checked residential: %s\n", username)
	if role == "manager" || role == "resident" {
		return true, nil
	} else {
		return false, nil
	}
}
