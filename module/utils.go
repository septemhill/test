package module

import (
	"context"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

var (
	number    = [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	lowerCase = [26]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	upperCase = [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
)

type RandomType int

const (
	RANDOM_NUM_ONLY = iota
	RANDOM_ALPHA_ONLY
	RANDOM_LCASE_WITH_NUM
	RANDOM_RCASE_WITH_NUM
	RANDOM_ALL
)

func numberCharArray() []byte {
	return number[:]
}

func alphaCharArray() []byte {
	b := append([]byte{}, lowerCase[:]...)
	return append(b, upperCase[:]...)
}

func lowerWithNumberCharArray() []byte {
	b := append([]byte{}, number[:]...)
	return append(b, lowerCase[:]...)
}

func upperWithNumberCharArray() []byte {
	b := append([]byte{}, number[:]...)
	return append(b, upperCase[:]...)
}

func allCharArray() []byte {
	b := append([]byte{}, number[:]...)
	b = append(b, lowerCase[:]...)
	return append(b, upperCase[:]...)
}

func GenRandomString(typ RandomType, length int) string {
	b := []byte{}

	switch typ {
	case RANDOM_NUM_ONLY:
		b = numberCharArray()
	case RANDOM_ALPHA_ONLY:
		b = alphaCharArray()
	case RANDOM_LCASE_WITH_NUM:
		b = lowerWithNumberCharArray()
	case RANDOM_RCASE_WITH_NUM:
		b = upperWithNumberCharArray()
	case RANDOM_ALL:
		b = allCharArray()
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
	return string([]byte{b[0], b[1], b[2], b[3], b[4], b[5]})
}

func generateRandomString() string {
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
