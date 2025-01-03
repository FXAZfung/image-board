package cmd

import (
	"github.com/FXAZfung/image-board/internal/db"
	"github.com/FXAZfung/image-board/internal/initialize"
	"github.com/FXAZfung/image-board/internal/initialize/data"
	"github.com/FXAZfung/image-board/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
)

func Init() {
	initialize.InitConfig()
	initialize.Log()
	initialize.InitDB()
	data.InitData()
}

func Release() {
	db.Close()
}

var pid = -1
var pidFile string

func initDaemon() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	_ = os.MkdirAll(filepath.Join(exPath, "daemon"), 0700)
	pidFile = filepath.Join(exPath, "daemon/pid")
	if utils.IsExist(pidFile) {
		bytes, err := os.ReadFile(pidFile)
		if err != nil {
			log.Fatal("failed to read pid file", err)
		}
		id, err := strconv.Atoi(string(bytes))
		if err != nil {
			log.Fatal("failed to parse pid data", err)
		}
		pid = id
	}
}
