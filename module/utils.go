package module

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
	ers "github.com/septemhill/test/errors"
)

func generateLink() string {
	var b = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D',
		'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
		'Y', 'Z'}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
	return string([]byte{b[0], b[1], b[2], b[3], b[4], b[5]})
}

func sessionTokenGenerate() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func txAction(ctx context.Context, db *db.DB, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = txFunc(tx)
	return err
}

func dbErrHandler(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Println(pgErr.Code)
		return ers.ErrParameter(err)
	}

	return ers.ErrUnknown(err)
}
