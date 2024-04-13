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

func CheckUser(user *model.User) (*model.User, error) {
	result, err := backend.PGBackend.ValidateUser(user.Username, user.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	fmt.Printf("Service chekced user: %s\n", user.Username)
	return result, nil
}
