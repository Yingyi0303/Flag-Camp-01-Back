package backend

import (
	"log"
	"propertyManagement/model"
)

func (backend *PostgresBackend) InsertReservation(username string, facilityId int, remark string, reservationDate string, startHour int, endHour int) (*model.Reservation, error) {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()

	var facilityName string
	err := tx.QueryRow("SELECT facility_name FROM facilities WHERE id = $1 FOR UPDATE", facilityId).Scan(&facilityName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	query := "SELECT COUNT(*) FROM reservations WHERE facility_name = $1 AND reservation_date = TO_DATE($2, 'YYYY-MM-DD') AND end_hour > $3 AND start_hour < $4"
	
	var count int
	err = tx.QueryRow(query, facilityName, reservationDate, startHour, endHour).Scan(&count)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if count > 0 {
		return nil, nil
	}

	query = "INSERT INTO reservations (username, facility_name, remark, reservation_date, start_hour, end_hour) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	var id int64
	err = tx.QueryRow(query, username, facilityName, remark, reservationDate, startHour, endHour).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := backend.SelectReservationById(int(id))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (backend *PostgresBackend) SelectAllReservationsNextWeek(startDate string) ([]model.Reservation, error) {
	rows, err := backend.db.Query("SELECT id, username, facility_name, remark, TO_CHAR(reservation_date, 'YYYY-MM-DD'), start_hour, end_hour FROM reservations WHERE reservation_date BETWEEN TO_DATE($1, 'YYYY-MM-DD') AND TO_DATE($1, 'YYYY-MM-DD') + INTERVAL '6 days' ORDER BY reservation_date, start_hour", startDate)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	reservations := []model.Reservation{}
	for rows.Next() {
		var reservation model.Reservation
		err := rows.Scan(&reservation.Id,
						 &reservation.Username,
						 &reservation.FacilityName,
						 &reservation.Remark,
						 &reservation.ReservationDate,
						 &reservation.StartHour,
						 &reservation.EndHour)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return reservations, nil
}

func (backend *PostgresBackend) SelectAllReservationsByFacilityName(facilityName, startDate string, endDate string) ([]model.Reservation, error) {
	query := "SELECT id, username, facility_name, remark, TO_CHAR(reservation_date, 'YYYY-MM-DD'), start_hour, end_hour FROM reservations WHERE facility_name = $1 AND reservation_date BETWEEN TO_DATE($2, 'YYYY-MM-DD') AND TO_DATE($3, 'YYYY-MM-DD') ORDER BY reservation_date, start_hour"

	rows, err := backend.db.Query(query, facilityName, startDate, endDate)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	reservations := []model.Reservation{}
	for rows.Next() {
		var reservation model.Reservation
		err := rows.Scan(&reservation.Id,
						 &reservation.Username,
						 &reservation.FacilityName,
						 &reservation.Remark,
						 &reservation.ReservationDate,
						 &reservation.StartHour,
						 &reservation.EndHour)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return reservations, nil
}

func (backend *PostgresBackend) SelectAllReservationsByUsername(username string) ([]model.Reservation, error) {
	rows, err := backend.db.Query("SELECT id, username, facility_name, remark, TO_CHAR(reservation_date, 'YYYY-MM-DD'), start_hour, end_hour FROM reservations WHERE username = $1 ORDER BY reservation_date DESC, start_hour DESC", username)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	reservations := []model.Reservation{}
	for rows.Next() {
		var reservation model.Reservation
		err := rows.Scan(&reservation.Id,
						 &reservation.Username,
						 &reservation.FacilityName,
						 &reservation.Remark,
						 &reservation.ReservationDate,
						 &reservation.StartHour,
						 &reservation.EndHour)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return reservations, nil
}

func (backend *PostgresBackend) SelectReservationById(id int) (*model.Reservation, error) {
	var reservation model.Reservation
	err := backend.db.QueryRow("SELECT id, username, facility_name, remark, TO_CHAR(reservation_date, 'YYYY-MM-DD'), start_hour, end_hour FROM reservations WHERE id = $1", id).
			Scan(&reservation.Id,
				 &reservation.Username,
				 &reservation.FacilityName,
				 &reservation.Remark,
				 &reservation.ReservationDate,
				 &reservation.StartHour,
				 &reservation.EndHour)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &reservation, nil
}

func (backend *PostgresBackend) DeleteReservation(username string, id int) error {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()

	query := "DELETE FROM reservations WHERE username = $1 AND id = $2"

	_, err := tx.Exec(query, username, id)
	if err != nil {
		log.Println(err)
		return err
	}
	
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
