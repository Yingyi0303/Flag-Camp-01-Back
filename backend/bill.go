package backend

import (
	"log"
	"propertyManagement/model"
	"time"
)

func (backend *PostgresBackend) InsertBill(username string, maintenance_id int, item string, amount int) (*model.Bill, error) {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()
	
	var isCompleted bool
	err := tx.QueryRow("SELECT completed from maintenances WHERE username = $1 AND id = $2 FOR UPDATE", username, maintenance_id).Scan(&isCompleted)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if isCompleted {
		return nil, nil
	}

	query := "INSERT INTO bills (username, maintenance_id, item, amount, bill_time) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	var id int64
	err = tx.QueryRow(query, username, maintenance_id, item, amount, formattedTime).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := backend.SelectBillById(int(id))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (backend *PostgresBackend) SelectAllBillsByUsername(username string) ([]model.Bill, error) {
	rows, err := backend.db.Query("SELECT id, username, maintenance_id, item, amount, TO_CHAR(bill_time, 'YYYY-MM-DD HH24:MI:SS') FROM bills WHERE username = $1 ORDER BY bill_time DESC", username)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	bills := []model.Bill{}
	for rows.Next() {
		var bill model.Bill
		err := rows.Scan(&bill.Id,
						 &bill.Username,
						 &bill.MaintenanceId,
						 &bill.Item,
						 &bill.Amount,
						 &bill.BillTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		bills = append(bills, bill)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return bills, nil
}

func (backend *PostgresBackend) SelectBillById(id int) (*model.Bill, error) {
	var bill model.Bill
	err := backend.db.QueryRow("SELECT id, username, maintenance_id, item, amount, TO_CHAR(bill_time, 'YYYY-MM-DD HH24:MI:SS') FROM bills WHERE id = $1", id).
		   Scan(&bill.Id,
				&bill.Username,
				&bill.MaintenanceId,
				&bill.Item,
				&bill.Amount,
				&bill.BillTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &bill, nil
}
