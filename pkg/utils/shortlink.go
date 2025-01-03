package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

func GenerateShortLink(input string) string {
	hash := sha1.Sum([]byte(input))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(encoded[:8], "=") // 截取前 8 位，去掉填充符
}
