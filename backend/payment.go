package backend

import (
	"log"
	"propertyManagement/model"
	"time"
)

func (backend *PostgresBackend) InsertPayment(username string, item string, amount int) (*model.Payment, error) {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()

	var balance int
	err := tx.QueryRow("SELECT amount from balances WHERE username = $1 FOR UPDATE", username).Scan(&balance)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if balance < amount {
		return nil, nil
	}

	query := "INSERT INTO payments (username, item, amount, payment_time) VALUES ($1, $2, $3, $4) RETURNING id"

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	var id int64
	err = tx.QueryRow(query, username, item, amount, formattedTime).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := backend.SelectPaymentById(int(id))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (backend *PostgresBackend) SelectAllPaymentsByUsername(username string) ([]model.Payment, error) {
	rows, err := backend.db.Query("SELECT id, username, item, amount, TO_CHAR(payment_time, 'YYYY-MM-DD HH24:MI:SS') FROM payments WHERE username = $1 ORDER BY payment_time DESC", username)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	payments := []model.Payment{}
	for rows.Next() {
		var payment model.Payment
		err := rows.Scan(&payment.Id,
						 &payment.Username,
						 &payment.Item,
						 &payment.Amount,
						 &payment.PaymentTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		payments = append(payments, payment)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return payments, nil
}

func (backend *PostgresBackend) SelectPaymentById(id int) (*model.Payment, error) {
	var payment model.Payment
	err := backend.db.QueryRow("SELECT id, username, item, amount, TO_CHAR(payment_time, 'YYYY-MM-DD HH24:MI:SS') FROM payments WHERE id = $1", id).
		   Scan(&payment.Id,
				&payment.Username,
				&payment.Item,
				&payment.Amount,
				&payment.PaymentTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &payment, nil
}
