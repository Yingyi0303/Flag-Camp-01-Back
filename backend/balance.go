package backend

import (
	"log"
)

func (backend *PostgresBackend) SelectBalanceByUsername (username string) (int, error) {
	var balance int
	err := backend.db.QueryRow("SELECT amount FROM balances WHERE username = $1", username).Scan(&balance)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return balance, err
}
