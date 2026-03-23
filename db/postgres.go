package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Connect() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		fmt.Println("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
