package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func AddMaintenance(maintenance *model.Maintenance) (*model.Maintenance, error) {
	maintenance, err := backend.PGBackend.InsertMaintenance(maintenance.Username, maintenance.Subject, maintenance.Content)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service added maintenance: %d\n", maintenance.Id)
	return maintenance, nil
}

func GetAllMaintenances(completed bool) ([]model.Maintenance, error) {
	maintenances, err := backend.PGBackend.SelectAllMaintenances(completed)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all maintenances completed: %t\n", completed)
	return maintenances, nil
}

func GetMyMaintenances(username string, completed bool) ([]model.Maintenance, error) {
	maintenances, err := backend.PGBackend.SelectAllMaintenancesByUsername(username, completed)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all maintenances by user: %s completed: %t\n", username, completed)
	return maintenances, nil
}

func SetMaintenance(maintenance *model.Maintenance) (*model.Maintenance, error) {
	success, err := backend.PGBackend.MaintenanceExists(maintenance.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !success {
		return nil, nil
	}
	result, err := backend.PGBackend.UpdateMaintenanceById(maintenance.Id, maintenance.Reply, maintenance.Completed)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	fmt.Printf("Service set maintenance: %d\n", maintenance.Id)
	return result, nil
}
