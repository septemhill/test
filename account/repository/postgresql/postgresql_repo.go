package postgresql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/account"
)

type accountRepository struct {
	*sqlx.DB
}

func NewAccountRepository() *accountRepository {
	//connInfo := fmt.Sprintf(`user=%s password=%s dbname=%s port=%s`,
	//	os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))
	//db := sqlx.MustConnect("pgx", connInfo)
	//
	//return &accountRepository{DB: db}
	return &accountRepository{}
}

func txAction(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
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

	err = fn(tx)
	return
}

func (repo *accountRepository) Create(ctx context.Context, acc *account.Account) (int, error) {
	var id int
	accExpr := `INSERT INTO accounts VALUES(DEFAULT, $1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := txAction(ctx, repo.DB, func(tx *sqlx.Tx) error {
		curr := time.Now().Truncate(time.Millisecond).UTC()
		if err := tx.GetContext(ctx, &id, accExpr, acc.Username, acc.Email, acc.Phone, acc.Password, "NORMAL", curr, curr); err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func (repo *accountRepository) GetInfo(ctx context.Context, id int) (*account.Account, error) {
	var acc account.Account
	expr := `SELECT email, phone FROM accounts WHERE id = $1`

	if err := txAction(ctx, repo.DB, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &acc, expr, id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &acc, nil
}

func (repo *accountRepository) UpdateInfo(ctx context.Context, acc *account.Account) (int, error) {
	var id int
	expr := `UPDATE accounts SET phone = $1 WHERE username = $2 RETURNING id`

	err := txAction(ctx, repo.DB, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, acc.Phone, acc.Username); err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func (repo *accountRepository) Delete(ctx context.Context, id int) (int, error) {
	var rid int
	expr := `DELETE FROM accounts WHERE id = $1 RETURNING id`
	err := txAction(ctx, repo.DB, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &rid, expr, id); err != nil {
			return err
		}

		return nil
	})

	return rid, err
}

func (repo *accountRepository) ChangePassword(ctx context.Context, id int, newPasswd string) (int, error) {
	var rid int
	expr := `UPDATE accounts SET password = $1 WHERE id = $2 RETURNING id`

	err := txAction(ctx, repo.DB, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &rid, expr, newPasswd, id); err != nil {
			return err
		}

		return nil
	})
	return rid, err
}

var _ account.AccountRepository = (*accountRepository)(nil)
