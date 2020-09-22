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
	var username string
	expr := `SELECT username FROM accounts WHERE email = $1 AND password = $2`
	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &username, expr, email, password); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	token := utils.GenerateRandomString(utils.RANDOM_HEX_ONLY, SESSION_TOKEN_LEN)

	h := map[string]interface{}{
		SESS_HSET_USERNAME: username,
		SESS_HSET_EMAIL:    email,
	}

	if _, err := r.HMSet(SessionTokenPrefix(token), h).Result(); err != nil {
		return "", nil
	}

	return token, nil
}

func Signup(ctx context.Context, d *db.DB, r *redis.Client, acc *Account) (string, error) {
	var id int
	accExpr := `INSERT INTO accounts VALUES(DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		curr := time.Now().Truncate(time.Millisecond).UTC()
		if err := tx.GetContext(ctx, &id, accExpr, acc.Username, acc.Email, acc.Phone, acc.Password, "NORMAL", curr, curr, false); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	code := utils.GenerateRandomString(utils.RANDOM_ALL, FORGET_PASSWD_LEN)
	key := SignupKeyPrefix(code)

	if _, err := r.Set(key, acc.Email, SignUpKeyTimeout).Result(); err != nil {
		return "", err
	}

	return code, nil
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

	var id string
	activateAccount := `UPDATE accounts SET active = true WHERE email = $1 RETURNING id`

	// TODO: postgres and redis should be in an atomic operation
	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, activateAccount, email); err != nil {
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
	expr := `UPDATE accounts SET password = $1 WHERE email = $2`

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
