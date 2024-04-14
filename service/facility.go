package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func GetAllFacilities() ([]model.Facility, error) {
	facilities, err := backend.PGBackend.SelectAllFacilities()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("Service fetched all facilities")
	return facilities, nil
}
