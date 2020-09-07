package module

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

func txAction(ctx context.Context, db *db.DB, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = txFunc(tx)
	return err
}
