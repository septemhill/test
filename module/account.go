package module

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
	"gopkg.in/guregu/null.v4"
)

type Account struct {
	ID       int         `db:"id"`
	Username string      `db:"username" json:"username"`
	Email    string      `db:"email" json:"email"`
	Password string      `db:"password" json:"password"`
	Phone    null.String `db:"phone" json:"phone"`
}

func CreateAccount(ctx context.Context, db *db.DB, acc Account) (err error) {
	accExpr := `INSERT INTO accounts VALUES(DEFAULT, $1, $2, $3)`
	accpriExpr := `INSERT INTO accounts_private VALUES (DEFAULT, $1, $2)`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, accExpr, acc.Username, acc.Email, acc.Phone); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, accpriExpr, acc.Email, acc.Password); err != nil {
			return err
		}

		return nil
	})
}

func GetAccountInfo(ctx context.Context, db *db.DB, acc *Account) error {
	expr := `SELECT email, phone FROM accounts WHERE username = $1`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, acc, expr, acc.Username); err != nil {
			return err
		}

		return nil
	})
}

func UpdateAccountInfo(ctx context.Context, db *db.DB, acc Account) error {
	expr := `UPDATE accounts SET phone = $1 WHERE username = $2`

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
	expr := `DELETE FROM accounts WHERE username = $1 AND email = $2`

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
