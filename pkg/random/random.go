package random

import (
	"crypto/rand"
	"github.com/google/uuid"
	"math/big"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func String(n int) string {
	b := make([]byte, n)
	letterLen := big.NewInt(int64(len(letterBytes)))
	for i := range b {
		idx, err := rand.Int(rand.Reader, letterLen)
		if err != nil {
			panic(err)
		}
		b[i] = letterBytes[idx.Int64()]
	}
	return string(b)
}

func SecretKey() string {
	return "image-board-" + uuid.NewString() + String(64)
}
