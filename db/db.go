package db

import (
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
	connInfo := `user=septemlee dbname=runtests sslmode=disable`
	db := sqlx.MustConnect("pgx", connInfo)
	return &DB{DB: db}
}
