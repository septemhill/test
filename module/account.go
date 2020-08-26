package module

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

type Account struct {
	ID       int    `db:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Phone    string `db:"phone" json:"phone"`
}

func CreateAccount(ctx context.Context, db *db.DB, acc Account) (err error) {
	expr := `INSER INTO account VALUES(DEFAULT, $1, $2, $3, $4)`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		res, err := tx.ExecContext(ctx, expr, acc.Username, acc.Email, acc.Phone, acc.Password)
		if err != nil {
			return err
		}

		count, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if count != 1 {
			return errors.New("insert affected row not exactly 1")
		}

		return nil
	})
}

func GetAccountInfo(ctx context.Context, db *db.DB, acc *Account) error {
	expr := `SELECT email, phone FROM account WHERE username = $1`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		row := tx.QueryRowxContext(ctx, expr, acc.Username)

		if err := row.StructScan(acc); err != nil {
			return err
		}

		return nil
	})
}

func UpdateAccountInfo(ctx context.Context, db *db.DB, acc Account) error {
	expr := `UPDATE account SET phone = $1 WHERE username = $2`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		res, err := tx.ExecContext(ctx, expr, acc.Phone, acc.Username)
		if err != nil {
			return err
		}

		count, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if count != 1 {
			return errors.New("update affected row not exactly 1")
		}

		return nil
	})
}

func DeleteAccount(ctx context.Context, db *db.DB, acc Account) error {
	expr := `DELETE FROM account WHERE username = $1 AND email = $2`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		res, err := tx.ExecContext(ctx, expr, acc.Username, acc.Email)
		if err != nil {
			return nil
		}

		count, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if count != 1 {
			return errors.New("delete affected row not exactly 1")
		}

		return nil
	})
}
