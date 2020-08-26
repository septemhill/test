package module

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

type SignupInfo struct {
	Username string
	Password string
	Email    string
	Phone    string
}

func Signup(ctx context.Context, db *db.DB, redis *redis.Client, info SignupInfo) (string, error) {
	expr := `INSERT INTO non_verified_accounts VALUES(DEFAULT, $1, $2, $3, $4)`

	err := txAction(ctx, db, func(tx *sqlx.Tx) error {
		res, err := tx.ExecContext(ctx, expr, info.Username, info.Password, info.Email, info.Phone)
		if err != nil {
			return err
		}

		count, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if count != 1 {
			return errors.New("signup affected row not exactly 1")
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	code := generateLink()
	key := SignupKeyPrefix(code)

	if _, err = redis.Set(key, info.Username, SignUpKeyTimeout).Result(); err != nil {
		return "", err
	}

	return key, nil
}

func VerifyUserRegistration(ctx context.Context, db *db.DB, redis *redis.Client, code string) error {
	key := SignupKeyPrefix(code)
	username, err := redis.Get(key).Result()
	if err != nil {
		return err
	}

	if username == "" {
		return errors.New("This link already expired")
	}

	insAccount := `INSERT INTO accounts SELECT nextval('accounts_id_seq'), username, email, phone FROM non_verified_accounts WHERE username = $1`
	insAccountPrivate := `INSERT INTO accounts_private SELECT nexeval('accounts_private_id_seq'), username, password FROM non_verified_accounts WHERE username = $1`
	delNonVerify := `DELETE FROM non_verified_accounts WHERE username = $1`

	// TODO: postgres and redis should be in an atomic operation
	err = txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, insAccount, username); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, insAccountPrivate, username); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, delNonVerify, username); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	_, err = redis.Del(key).Result()
	if err != nil {
		return err
	}

	return nil
}
