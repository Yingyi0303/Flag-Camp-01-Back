package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
	"time"
)

func AddReservation(reservation *model.Reservation) (*model.Reservation, error) {
	facilityId, err := backend.PGBackend.SelectIdByFacilityName(reservation.FacilityName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result, err := backend.PGBackend.InsertReservation(reservation.Username,
													   facilityId,
													   reservation.Remark,
													   reservation.ReservationDate,
													   reservation.StartHour,
													   reservation.EndHour)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	fmt.Printf("Service added reservation: %d\n", result.Id)
	return result, nil
}

func GetAllReservations() ([]model.Reservation, error) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	reservations, err := backend.PGBackend.SelectAllReservationsNextWeek(currentDate)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("Service fetched all reservations next week")
	return reservations, nil
}

func GetFacilityReservations(request *model.ReservationRequest) ([]model.Reservation, error) {
	reservations, err := backend.PGBackend.SelectAllReservationsByFacilityName(request.FacilityName, request.StartDate, request.EndDate)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all reservations by facility: %s\n", request.FacilityName)
	return reservations, nil
}

func GetMyReservations(username string) ([]model.Reservation, error) {
	reservations, err := backend.PGBackend.SelectAllReservationsByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all reservations by user: %s\n", username)
	return reservations, nil
}

func RemoveReservation(username string, id int) error {
	err := backend.PGBackend.DeleteReservation(username, id)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("Service removed discussion: %d\n", id)
	return nil
}
