package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func AddPayment(payment *model.Payment) (*model.Payment, error) {
	result, err := backend.PGBackend.InsertPayment(payment.Username, payment.Item, payment.Amount)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	fmt.Printf("Service added payment: %d\n", result.Id)
	return result, nil
}

func GetMyPayments(username string) ([]model.Payment, error) {
	payments, err := backend.PGBackend.SelectAllPaymentsByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all payments by user: %s\n", username)
	return payments, nil
}
