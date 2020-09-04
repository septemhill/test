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

func Login(ctx context.Context, db *db.DB, redis *redis.Client, username, password string) (string, error) {
	expr := `SELECT COUNT(*) FROM accounts_private WHERE username = $1 AND password = $2`

	if err := txAction(ctx, db, func(tx *sqlx.Tx) error {
		res := tx.QueryRowxContext(ctx, expr, username, password)

		cnt := 0
		if err := res.Scan(&cnt); err != nil {
			return dbErrHandler(err)
		}

		if cnt < 1 {
			return dbErrHandler(errors.New("Invalide username/password"))
		}

		return nil
	}); err != nil {
		return "", err
	}

	token := sessionTokenGenerate()
	redis.Set(token, username, time.Hour*1)

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

	if _, err := redis.Set(key, info.Username, SignUpKeyTimeout).Result(); err != nil {
		return "", err
	}

	return key, nil
}

func VerifyUserRegistration(ctx context.Context, db *db.DB, redis *redis.Client, code string) error {
	key := SignupKeyPrefix(code)
	username, err := redis.Get(key).Result()
	if err != nil {
		return dbErrHandler(err)
	}

	if username == "" {
		return dbErrHandler(errors.New("This link already expired"))
	}

	insAccount := `INSERT INTO accounts SELECT nextval('accounts_id_seq'), username, email, phone FROM non_verified_accounts WHERE username = $1`
	insAccountPrivate := `INSERT INTO accounts_private SELECT nexeval('accounts_private_id_seq'), username, password FROM non_verified_accounts WHERE username = $1`
	delNonVerify := `DELETE FROM non_verified_accounts WHERE username = $1`

	// TODO: postgres and redis should be in an atomic operation
	if err = txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, insAccount, username); err != nil {
			return dbErrHandler(err)
		}

		if _, err := tx.ExecContext(ctx, insAccountPrivate, username); err != nil {
			return dbErrHandler(err)
		}

		if _, err := tx.ExecContext(ctx, delNonVerify, username); err != nil {
			return dbErrHandler(err)
		}

		return nil
	}); err != nil {
		return dbErrHandler(err)
	}

	_, err = redis.Del(key).Result()
	if err != nil {
		return dbErrHandler(err)
	}

	return nil
}
