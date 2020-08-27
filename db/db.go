package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func OpenDB() *DB {
	connInfo := `user=septemlee dbname=septemlee sslmode=disable`
	db := sqlx.MustConnect("pgx", connInfo)
	return &DB{DB: db}
}

func (db *DB) executeSQL(ctx context.Context, expr string, args ...interface{}) (sql.Result, error) {
	return db.ExecContext(ctx, expr, args...)
}

func (db *DB) Insert(ctx context.Context, expr string, args ...interface{}) (sql.Result, error) {
	return db.executeSQL(ctx, expr, args...)
}

func (db *DB) Update(ctx context.Context, expr string, args ...interface{}) (sql.Result, error) {
	return db.executeSQL(ctx, expr, args...)
}

func (db *DB) Delete(ctx context.Context, expr string, args ...interface{}) (sql.Result, error) {
	return db.executeSQL(ctx, expr, args...)
}
