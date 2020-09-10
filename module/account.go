package module

import (
	"context"

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

func CreateAccount(ctx context.Context, d *db.DB, acc Account) (err error) {
	accExpr := `INSERT INTO accounts VALUES(DEFAULT, $1, $2, $3) RETURNING id`
	accPriExpr := `INSERT INTO accounts_private VALUES (DEFAULT, $1, $2) RETURNING id`

	return txAction(ctx, d, func(tx *sqlx.Tx) error {
		var id int
		if err := tx.GetContext(ctx, &id, accExpr, acc.Username, acc.Email, acc.Phone); err != nil {
			return err
		}

		if err := tx.GetContext(ctx, &id, accPriExpr, acc.Email, acc.Password); err != nil {
			return err
		}

		return nil
	})
}

func GetAccountInfo(ctx context.Context, d *db.DB, acc *Account) error {
	expr := `SELECT email, phone FROM accounts WHERE username = $1`

	return txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, acc, expr, acc.Username); err != nil {
			return err
		}

		return nil
	})
}

func UpdateAccountInfo(ctx context.Context, d *db.DB, acc Account) error {
	expr := `UPDATE accounts SET phone = $1 WHERE username = $2 RETURNING id`

	return txAction(ctx, d, func(tx *sqlx.Tx) error {
		var id int
		if err := tx.GetContext(ctx, &id, expr, acc.Phone, acc.Username); err != nil {
			return err
		}

		return nil
	})
}

func DeleteAccount(ctx context.Context, d *db.DB, acc Account) error {
	accExpr := `DELETE FROM accounts WHERE username = $1 AND email = $2 RETURNING id`
	accPriExpr := `DELETE FROM accounts_private WHERE email = $1 RETURNING id`

	return txAction(ctx, d, func(tx *sqlx.Tx) error {
		var id int

		if err := tx.GetContext(ctx, &id, accPriExpr, acc.Email); err != nil {
			return err
		}

		if err := tx.GetContext(ctx, &id, accExpr, acc.Username, acc.Email); err != nil {
			return err
		}

		return nil
	})
}
