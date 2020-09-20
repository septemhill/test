package module

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
	"gopkg.in/guregu/null.v4"
)

type Account struct {
	ID       int         `db:"id"`
	Username string      `db:"username" json:"username" uri:"user"`
	Email    string      `db:"email" json:"email"`
	Password string      `db:"password" json:"password"`
	Phone    null.String `db:"phone" json:"phone"`
	CreateAt time.Time   `db:"create_at" json:"createAt"`
	UpdateAt time.Time   `db:"update_at" json:"updateAt"`
}

func CreateAccount(ctx context.Context, d *db.DB, acc *Account) (int, error) {
	var id int
	accExpr := `INSERT INTO accounts VALUES(DEFAULT, $1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		curr := time.Now().Truncate(time.Millisecond).UTC()
		if err := tx.GetContext(ctx, &id, accExpr, acc.Username, acc.Email, acc.Phone, acc.Password, "NORMAL", curr, curr); err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func GetAccountInfo(ctx context.Context, d *db.DB, acc *Account) (*Account, error) {
	expr := `SELECT email, phone FROM accounts WHERE username = $1`

	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, acc, expr, acc.Username); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return acc, nil
}

func UpdateAccountInfo(ctx context.Context, d *db.DB, acc *Account) (int, error) {
	var id int
	expr := `UPDATE accounts SET phone = $1 WHERE username = $2 RETURNING id`

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, acc.Phone, acc.Username); err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func DeleteAccount(ctx context.Context, d *db.DB, acc *Account) (int, error) {
	var id int
	expr := `DELETE FROM accounts WHERE username = $1 AND email = $2 RETURNING id`
	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, acc.Username, acc.Email); err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func ChangePassword(ctx context.Context, d *db.DB, email, newPassword string) (int, error) {
	var id int
	expr := `UPDATE accounts SET password = $1 WHERE email = $2 RETURNING id`

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, newPassword, email); err != nil {
			return err
		}

		return nil
	})

	return id, err
}
