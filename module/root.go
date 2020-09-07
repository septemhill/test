package module

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
	"gopkg.in/guregu/null.v4"
)

type SignupInfo struct {
	Username string      `db:"username" json:"username"`
	Password string      `db:"password" json:"password"`
	Email    string      `db:"email" json:"email"`
	Phone    null.String `db:"phone" json:"phone"`
}

func Login(ctx context.Context, db *db.DB, redis *redis.Client, email, password string) (string, error) {
	expr := `SELECT COUNT(*) FROM accounts_private WHERE email = $1 AND password = $2`

	if err := txAction(ctx, db, func(tx *sqlx.Tx) error {
		cnt := 0
		if err := tx.GetContext(ctx, &cnt, expr, email, password); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	token := sessionTokenGenerate()

	if _, err := redis.Set(token, email, time.Hour*1).Result(); err != nil {
		return "", err
	}

	return token, nil
}

func Signup(ctx context.Context, db *db.DB, redis *redis.Client, info SignupInfo) (string, error) {
	expr := `INSERT INTO non_verified_accounts VALUES(DEFAULT, $1, $2, $3, $4)`

	if err := txAction(ctx, db, func(tx *sqlx.Tx) error {
		res, err := tx.ExecContext(ctx, expr, info.Username, info.Password, info.Email, info.Phone)
		if err != nil {
			return dbErrHandler(err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			return dbErrHandler(err)
		}

		if count != 1 {
			return dbErrHandler(errors.New("signup affected row not exactly 1"))
		}

		return nil
	}); err != nil {
		return "", err
	}

	code := generateLink()
	key := SignupKeyPrefix(code)

	if _, err := redis.Set(key, info.Email, SignUpKeyTimeout).Result(); err != nil {
		return "", err
	}

	return key, nil
}

func VerifyUserRegistration(ctx context.Context, db *db.DB, redis *redis.Client, code string) error {
	key := SignupKeyPrefix(code)
	email, err := redis.Get(key).Result()
	if err != nil {
		return dbErrHandler(err)
	}

	if email == "" {
		return errors.New("This link already expired")
	}

	insAccount := `INSERT INTO accounts SELECT nextval('accounts_id_seq'), username, email, phone FROM non_verified_accounts WHERE username = $1`
	insAccountPrivate := `INSERT INTO accounts_private SELECT nexeval('accounts_private_id_seq'), email, password FROM non_verified_accounts WHERE email = $1`
	delNonVerify := `DELETE FROM non_verified_accounts WHERE username = $1`

	// TODO: postgres and redis should be in an atomic operation
	if err = txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, insAccount, email); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, insAccountPrivate, email); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, delNonVerify, email); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return dbErrHandler(err)
	}

	if _, err := redis.Del(key).Result(); err != nil {
		return err
	}

	return nil
}

func ForgetPassword(ctx context.Context, db *db.DB, redis *redis.Client, email string) (string, error) {
	expr := `SELECT * FROM accounts WHERE email = $1`

	acc := Account{}

	// 1. Check email exist
	if err := txAction(ctx, db, func(tx *sqlx.Tx) error {
		if err := db.GetContext(ctx, &acc, expr, email); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	// 2. Generate hash code
	code := generateLink()

	// 3. Set hash code in redis
	if _, err := redis.Set(ForgetPasswordKeyPrefix(code), email, ForgetPasswdKeyTimeout).Result(); err != nil {
		return "", nil
	}

	return code, nil
}

func ResetPassword(ctx context.Context, db *db.DB, redis *redis.Client, code, password string) error {
	expr := `UPDATE accounts_private SET password = $1 WHERE email = $2`

	email, err := redis.Get(code).Result()
	if err != nil {
		return err
	}

	if email == "" {
		return errors.New("This link already expired")
	}

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, password, email); err != nil {
			return err
		}
		return nil
	})
}
