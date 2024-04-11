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