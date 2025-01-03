package common

import (
	"github.com/FXAZfung/go-cache"
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

var SecretKey []byte

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var validTokenCache = cache.NewMemCache[bool]()

func GenerateToken(user *model.User) (tokenString string, err error) {
	claim := UserClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(conf.Conf.TokenExpiresIn) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	validTokenCache.Set(tokenString, true)
	return tokenString, err
}

func ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if IsTokenInvalidated(tokenString) {
		return nil, errors.New("token is invalidated")
	}
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

func InvalidateToken(tokenString string) error {
	if tokenString == "" {
		return nil // don't invalidate empty guest token
	}
	validTokenCache.Del(tokenString)
	return nil
}

func IsTokenInvalidated(tokenString string) bool {
	_, ok := validTokenCache.Get(tokenString)
	return !ok
}
