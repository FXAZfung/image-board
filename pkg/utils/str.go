package utils

import (
	"encoding/base64"
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/errs"
	"strings"
)

func MappingName(name string) string {
	for k, v := range conf.FilenameCharMap {
		name = strings.ReplaceAll(name, k, v)
	}
	return name
}

var DEC = map[string]string{
	"-": "+",
	"_": "/",
	".": "=",
}

func SafeAtob(data string) (string, error) {
	for k, v := range DEC {
		data = strings.ReplaceAll(data, k, v)
	}
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// GetNoneEmpty returns the first non-empty string, return empty if all empty
func GetNoneEmpty(strArr ...string) string {
	for _, s := range strArr {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

// TrimSpace 去除字符串两端的空格，返回去除空格后的字符串，同时不能出现完全是空格的字符串，如果处理之后长度为0则抛出错误
func TrimSpace(str string) (string, error) {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return "", errs.ErrEmptyString
	}
	return str, nil
}
