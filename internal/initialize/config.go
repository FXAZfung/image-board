package initialize

import (
	"github.com/FXAZfung/image-board/cmd/flags"
	"github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func InitConfig() {
	// 如果是第一次启动，初始化默认配置
	configPath := filepath.Join(flags.DataDir, "config.json")
	imagesPath := filepath.Join(flags.DataDir, "images")
	exist := utils.IsExist(configPath)
	if !exist {
		log.Infof("config file not exists, creating default config file")
		config.Conf = config.DefaultConfig()
		_, err := utils.CreateNestedFile(configPath)
		if err != nil {
			log.Fatalf("failed to create config file: %+v", err)
		}
		err = utils.CreateNestedDirectory(imagesPath)
		if err != nil {
			log.Fatalf("failed to create images directory: %+v", err)
		}
		if !utils.WriteJsonToFile(configPath, config.Conf) {
			log.Fatalf("failed to create default config file")
		}
	} else {
		configBytes, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("reading config file error: %+v", err)
		}
		config.Conf = config.DefaultConfig()
		err = utils.Json.Unmarshal(configBytes, config.Conf)
		if err != nil {
			log.Fatalf("load config error: %+v", err)
		}
		// update config.json struct
		confBody, err := utils.Json.MarshalIndent(config.Conf, "", "  ")
		if err != nil {
			log.Fatalf("marshal config error: %+v", err)
		}
		err = os.WriteFile(configPath, confBody, 0o777)
		if err != nil {
			log.Fatalf("update config struct error: %+v", err)
		}
	}
}
