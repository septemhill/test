package module

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/utils"
	"gopkg.in/guregu/null.v4"
)

type SignupInfo struct {
	Username string      `db:"username" json:"username"`
	Password string      `db:"password" json:"password"`
	Email    string      `db:"email" json:"email"`
	Phone    null.String `db:"phone" json:"phone"`
}

func Login(ctx context.Context, d *db.DB, r *redis.Client, email, password string) (string, error) {
	expr := `SELECT COUNT(*) FROM accounts_private WHERE email = $1 AND password = $2`

	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		cnt := 0
		if err := tx.GetContext(ctx, &cnt, expr, email, password); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	token := utils.GenerateRandomString(utils.RANDOM_HEX_ONLY, SESSION_TOKEN_LEN)

	if _, err := r.Set(token, email, time.Hour*1).Result(); err != nil {
		return "", err
	}

	return token, nil
}

func Signup(ctx context.Context, d *db.DB, r *redis.Client, info SignupInfo) (string, error) {
	expr := `INSERT INTO non_verified_accounts VALUES(DEFAULT, $1, $2, $3, $4) RETURNING id`

	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		var id int
		if err := tx.GetContext(ctx, &id, expr, info.Username, info.Password, info.Email, info.Phone); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	code := utils.GenerateRandomString(utils.RANDOM_ALL, FORGET_PASSWD_LEN)
	key := SignupKeyPrefix(code)

	if _, err := r.Set(key, info.Email, SignUpKeyTimeout).Result(); err != nil {
		return "", err
	}

	return key, nil
}

func VerifyUserRegistration(ctx context.Context, d *db.DB, r *redis.Client, code string) error {
	key := SignupKeyPrefix(code)
	email, err := r.Get(key).Result()
	if err != nil {
		return err
	}

	if email == "" {
		return errors.New("link already expired")
	}

	insAccount := `
		INSERT INTO accounts SELECT nextval('accounts_id_seq'), username, email, phone FROM non_verified_accounts WHERE username = $1
	`
	insAccountPrivate := `
		INSERT INTO accounts_private SELECT nexeval('accounts_private_id_seq'), email, password FROM non_verified_accounts WHERE email = $1
	`
	delNonVerify := `
		DELETE FROM non_verified_accounts WHERE username = $1
	`

	// TODO: postgres and redis should be in an atomic operation
	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
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
		return err
	}

	if _, err := r.Del(key).Result(); err != nil {
		return err
	}

	return nil
}

func ForgetPassword(ctx context.Context, d *db.DB, r *redis.Client, email string) (string, error) {
	expr := `SELECT * FROM accounts WHERE email = $1`

	acc := Account{}

	// 1. Check email exist
	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &acc, expr, email); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	// 2. Generate hash code
	code := utils.GenerateRandomString(utils.RANDOM_ALL, SIGNUP_TOKEN_LEN)

	// 3. Set hash code in redis
	if _, err := r.Set(ForgetPasswordKeyPrefix(code), email, ForgetPasswdKeyTimeout).Result(); err != nil {
		return "", nil
	}

	return code, nil
}

func ResetPassword(ctx context.Context, d *db.DB, r *redis.Client, code, password string) error {
	expr := `UPDATE accounts_private SET password = $1 WHERE email = $2`

	email, err := r.Get(code).Result()
	if err != nil {
		return err
	}

	if email == "" {
		return errors.New("link already expired")
	}

	return txAction(ctx, d, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, password, email); err != nil {
			return err
		}
		return nil
	})
}
