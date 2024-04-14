package backend

import (
	"log"
	"propertyManagement/model"
)

func (backend *PostgresBackend) SelectAllFacilities() ([]model.Facility, error) {
	rows, err := backend.db.Query("SELECT id, facility_name, description FROM facilities ORDER BY id")
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	facilities := []model.Facility{}
	for rows.Next() {
		var facility model.Facility
		err := rows.Scan(&facility.Id,
						 &facility.FacilityName,
						 &facility.Description)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		facilities = append(facilities, facility)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return facilities, nil
}

func (backend *PostgresBackend) SelectIdByFacilityName(facilityName string) (int, error) {
	var id int64
	err := backend.db.QueryRow("SELECT id FROM facilities WHERE facility_name = $1", facilityName).Scan(&id)
	if err != nil {
		log.Println(err)
        return -1, err
    }
	return int(id), nil
}
