package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var PGBackend *PostgresBackend

type PostgresBackend struct {
	db *sql.DB
}

func Init() {
	// connect to database
	connectionString := "host=localhost port=5432 user=postgres password=secret dbname=property_management sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Connected to database successfully")

	// initialize tables
	script, _ := os.ReadFile("init.sql")
	statements := strings.Split(string(script), ";")
	for _, statement := range statements {
		query := strings.TrimSpace(statement)
		if query != "" {
			_, err = db.Exec(query)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}

	fmt.Println("Initialzed tables successfully")

	PGBackend = &PostgresBackend{db: db}
}

func Close() {
	PGBackend.db.Close()
}

func (backend *PostgresBackend) InsertUser(username, password, role string) error {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()

	query := "INSERT INTO users (username, password, role) VALUES ($1, $2, $3)"
	_, err := tx.Exec(query, username, password, role)
	if err != nil {
		log.Fatal(err)
		return err
	}
	
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (backend *PostgresBackend) UserExists(username string) (bool, error) {
	var count int
	err := backend.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		log.Fatal(err)
        return false, err
    }
	return count > 0, nil
}

func (backend *PostgresBackend) ValidateUser(username, password string) (bool, error) {
	var count int
	err := backend.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2", username, password).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return count > 0, nil
}