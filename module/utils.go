package module

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

// func readUncommittedTxAction(ctx context.Context, d *db.DB, txFunc func(*sqlx.Tx) error) (err error) {
// 	return txAction(ctx, d, &sql.TxOptions{Isolation: sql.LevelReadUncommitted}, txFunc)
// }

func readCommittedTxAction(ctx context.Context, d *db.DB, txFunc func(*sqlx.Tx) error) (err error) {
	return txAction(ctx, d, &sql.TxOptions{Isolation: sql.LevelReadCommitted}, txFunc)
}

// func repeatableReadTxAction(ctx context.Context, d *db.DB, txFunc func(*sqlx.Tx) error) (err error) {
// 	return txAction(ctx, d, &sql.TxOptions{Isolation: sql.LevelRepeatableRead}, txFunc)
// }

// func serializableTxAction(ctx context.Context, d *db.DB, txFunc func(*sqlx.Tx) error) (err error) {
// 	return txAction(ctx, d, &sql.TxOptions{Isolation: sql.LevelSerializable}, txFunc)
// }

func txAction(ctx context.Context, d *db.DB, opts *sql.TxOptions, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := d.BeginTxx(ctx, opts)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = txFunc(tx)
	return err
}
