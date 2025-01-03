package sign

import (
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/setting"
	"github.com/FXAZfung/image-board/pkg/sign"
	"sync"
	"time"
)

var once sync.Once
var instance sign.Sign

func Sign(data string) string {
	expire := setting.GetInt(conf.LinkExpiration, 0)
	if expire == 0 {
		return NotExpired(data)
	} else {
		return WithDuration(data, time.Duration(expire)*time.Hour)
	}
}

func WithDuration(data string, d time.Duration) string {
	once.Do(Instance)
	return instance.Sign(data, time.Now().Add(d).Unix())
}

func NotExpired(data string) string {
	once.Do(Instance)
	return instance.Sign(data, 0)
}

func Verify(data string, sign string) error {
	once.Do(Instance)
	return instance.Verify(data, sign)
}

func Instance() {
	instance = sign.NewHMACSign([]byte(setting.GetStr(conf.Token)))
}
