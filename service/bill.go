package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func AddBill(bill *model.Bill) (*model.Bill, error) {
	result, err := backend.PGBackend.InsertBill(bill.Username, bill.MaintenanceId, bill.Item, bill.Amount)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	fmt.Printf("Service added bill: %d\n", result.Id)
	return result, nil
}

func GetMyBills(username string) ([]model.Bill, error) {
	bills, err := backend.PGBackend.SelectAllBillsByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all bills by user: %s\n", username)
	return bills, nil
}
