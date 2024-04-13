package backend

import (
	"log"
	"propertyManagement/model"
)

func (backend *PostgresBackend) InsertUser(username, password, role string) error {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()

	query := "INSERT INTO users (username, password, role) VALUES ($1, $2, $3)"
	_, err := tx.Exec(query, username, password, role)
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

func (backend *PostgresBackend) UserExists(username string) (bool, error) {
	var count int
	err := backend.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		log.Println(err)
        return false, err
    }
	return count > 0, nil
}

func (backend *PostgresBackend) ValidateUser(username, password string) (*model.User, error) {
	var user model.User
	err := backend.db.QueryRow("SELECT username, role FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user.Username, &user.Role)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
