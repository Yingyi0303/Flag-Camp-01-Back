package backend

import (
	"log"
	"propertyManagement/model"
	"strings"
	"time"
)

func (backend *PostgresBackend) InsertMaintenance(username, subject, content string) (*model.Maintenance, error) {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()
	
	query := "INSERT INTO maintenances (username, subject, content, last_update_time) VALUES ($1, $2, $3, $4) RETURNING id"

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	var id int
	err := tx.QueryRow(query, username, subject, content, formattedTime).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	return backend.SelectMaintenanceById(int(id))
}

func (backend *PostgresBackend) SelectAllMaintenances(completed bool) ([]model.Maintenance, error) {
	rows, err := backend.db.Query("SELECT id, username, subject, content, reply, completed, TO_CHAR(last_update_time, 'YYYY-MM-DD HH24:MI:SS') FROM maintenances WHERE completed = $1 ORDER BY last_update_time DESC", completed)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	maintenances := []model.Maintenance{}
	for rows.Next() {
		var maintenance model.Maintenance
		err := rows.Scan(&maintenance.Id,
						 &maintenance.Username,
						 &maintenance.Subject,
						 &maintenance.Content,
						 &maintenance.Reply,
						 &maintenance.Completed,
						 &maintenance.LastUpdateTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		maintenances = append(maintenances, maintenance)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return maintenances, nil
}

func (backend *PostgresBackend) SelectAllMaintenancesByUsername(username string, completed bool) ([]model.Maintenance, error) {
	rows, err := backend.db.Query("SELECT id, username, subject, content, reply, completed, TO_CHAR(last_update_time, 'YYYY-MM-DD HH24:MI:SS') FROM maintenances WHERE username = $1 AND completed = $2 ORDER BY last_update_time DESC", username, completed)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	maintenances := []model.Maintenance{}
	for rows.Next() {
		var maintenance model.Maintenance
		err := rows.Scan(&maintenance.Id,
						 &maintenance.Username,
						 &maintenance.Subject,
						 &maintenance.Content,
						 &maintenance.Reply,
						 &maintenance.Completed,
						 &maintenance.LastUpdateTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		maintenances = append(maintenances, maintenance)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return maintenances, nil
}

func (backend *PostgresBackend) UpdateMaintenanceById(id int, reply string, completed bool) (*model.Maintenance, error) {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()
	
	var isCompleted bool
	err := tx.QueryRow("SELECT completed from maintenances WHERE id = $1 FOR UPDATE", id).Scan(&isCompleted)
	if (err != nil) {
		log.Println(err)
		return nil, err
	}
	if isCompleted {
		return nil, nil
	}

	query := "UPDATE maintenances SET reply = reply || $1, completed = $2, last_update_time = $3 WHERE id = $4"

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	if len(reply) > 0 {
		reply = strings.TrimSpace(reply) + "\n"
	}

	_, err = tx.Exec(query, reply, completed, formattedTime, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	return backend.SelectMaintenanceById(int(id))
}

func (backend *PostgresBackend) SelectMaintenanceById(id int) (*model.Maintenance, error) {
	var maintenance model.Maintenance
	err := backend.db.QueryRow("SELECT id, username, subject, content, reply, completed, TO_CHAR(last_update_time, 'YYYY-MM-DD HH24:MI:SS') FROM maintenances WHERE id = $1", id).
		   Scan(&maintenance.Id,
				&maintenance.Username,
				&maintenance.Subject,
				&maintenance.Content,
				&maintenance.Reply,
				&maintenance.Completed,
				&maintenance.LastUpdateTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &maintenance, nil
}

func (backend *PostgresBackend) MaintenanceExists(id int) (bool, error) {
	var count int
	err := backend.db.QueryRow("SELECT COUNT(*) FROM maintenances WHERE id = $1", id).Scan(&count)
	if err != nil {
		log.Println(err)
        return false, err
    }
	return count > 0, nil
}
