package db

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func Connect() *DB {
	connInfo := `user=septemlee dbname=postgres sslmode=disable`
	db := sqlx.MustConnect("pgx", connInfo)
	return &DB{DB: db}
}

func ConnectToTest() *DB {
	connInfo := fmt.Sprintf(`user=%s password=%s`, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"))
	db := sqlx.MustConnect("pgx", connInfo)
	return &DB{DB: db}
}
