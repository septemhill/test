package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func Connect() *DB {
	connInfo := fmt.Sprintf(`user=%s password=%s dbname=%s port=%s`,
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))
	db := sqlx.MustConnect("pgx", connInfo)
	return &DB{DB: db}
}

func ConnectToTest() *DB {
	connInfo := fmt.Sprintf(`user=%s password=%s port=%s`,
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_PORT"))
	db := sqlx.MustConnect("pgx", connInfo)
	return &DB{DB: db}
}
