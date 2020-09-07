package utils

import (
	"math/rand"
	"time"
)

var (
	hex       = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
	digit     = [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	lowerCase = [26]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	upperCase = [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
)

type RandomType int

const (
	RANDOM_HEX_ONLY = iota
	RANDOM_DIGIT_ONLY
	RANDOM_ALPHA_ONLY
	RANDOM_LCASE_WITH_NUM
	RANDOM_RCASE_WITH_NUM
	RANDOM_ALL
)

func hexCharArray() []byte {
	return hex[:]
}

func numberCharArray() []byte {
	return digit[:]
}

func alphaCharArray() []byte {
	b := append([]byte{}, lowerCase[:]...)
	return append(b, upperCase[:]...)
}

func lowerWithNumberCharArray() []byte {
	b := append([]byte{}, digit[:]...)
	return append(b, lowerCase[:]...)
}

func upperWithNumberCharArray() []byte {
	b := append([]byte{}, digit[:]...)
	return append(b, upperCase[:]...)
}

func allCharArray() []byte {
	b := append([]byte{}, digit[:]...)
	b = append(b, lowerCase[:]...)
	return append(b, upperCase[:]...)
}

func GenerateRandomString(typ RandomType, length int) (str string) {
	b := []byte{}

	switch typ {
	case RANDOM_HEX_ONLY:
		b = hexCharArray()
	case RANDOM_DIGIT_ONLY:
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

	for i := 0; i < length; i++ {
		rand.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
		str += string(b[0])
	}
	return
}
