package random

import (
	"crypto/rand"
	"github.com/google/uuid"
	"math/big"
	"strings"
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

func UUID() string {
	return uuid.NewString()
}

// RandomizeFileName 获取原文件名后进行随机化处理保留文件格式
func RandomizeFileName(fileName string) string {
	// 获取文件名后缀
	suffix := fileName[strings.LastIndex(fileName, "."):]
	// 获取文件名前缀
	prefix := fileName[:strings.LastIndex(fileName, ".")]
	// 随机化文件名
	return prefix + "-" + String(8) + suffix
}

func SecretKey() string {
	return "IM-" + uuid.NewString() + String(64)
}
