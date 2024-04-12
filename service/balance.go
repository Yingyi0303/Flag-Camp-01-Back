package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func GetMyBalance(username string) (*model.BalanceDto, error) {
	balance, err := backend.PGBackend.SelectBalanceByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	bills, err := backend.PGBackend.SelectAllBillsByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	payments, err := backend.PGBackend.SelectAllPaymentsByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	balanceDto := model.BalanceDto {
		Balance: balance,
		Bills: bills,
		Payments: payments,
	}
	fmt.Printf("Service fetched balance by user: %s\n", username)
	return &balanceDto, nil
}
