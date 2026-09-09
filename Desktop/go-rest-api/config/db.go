package config


import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := "host=localhost port=5433 user=postgres password=12345 dbname=expensetracker sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Baza bağlantı xətası: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Baza cavab vermir: %v", err)
	}

	fmt.Println("PostgreSQL bazasına uğurla qoşuldu!")
}