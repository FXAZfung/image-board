package config

import (
	"github.com/FXAZfung/image-board/cmd/flags"
	"github.com/FXAZfung/image-board/pkg/random"
	"path/filepath"
)

type Database struct {
	Type        string `json:"type" env:"TYPE"`
	Host        string `json:"host" env:"HOST"`
	Port        int    `json:"port" env:"PORT"`
	User        string `json:"user" env:"USER"`
	Password    string `json:"password" env:"PASS"`
	Name        string `json:"name" env:"NAME"`
	DBFile      string `json:"db_file" env:"FILE"`
	TablePrefix string `json:"table_prefix" env:"TABLE_PREFIX"`
	SSLMode     string `json:"ssl_mode" env:"SSL_MODE"`
	DSN         string `json:"dsn" env:"DSN"`
}

type Cors struct {
	AllowOrigins []string `json:"allow_origins" env:"ALLOW_ORIGINS"`
	AllowMethods []string `json:"allow_methods" env:"ALLOW_METHODS"`
	AllowHeaders []string `json:"allow_headers" env:"ALLOW_HEADERS"`
}

type LogConfig struct {
	Enable     bool   `json:"enable" env:"LOG_ENABLE"`
	Name       string `json:"name" env:"LOG_NAME"`
	MaxSize    int    `json:"max_size" env:"MAX_SIZE"`
	MaxBackups int    `json:"max_backups" env:"MAX_BACKUPS"`
	MaxAge     int    `json:"max_age" env:"MAX_AGE"`
	Compress   bool   `json:"compress" env:"COMPRESS"`
}

type Scheme struct {
	Address      string `json:"address" env:"ADDR"`
	HttpPort     int    `json:"http_port" env:"HTTP_PORT"`
	HttpsPort    int    `json:"https_port" env:"HTTPS_PORT"`
	ForceHttps   bool   `json:"force_https" env:"FORCE_HTTPS"`
	CertFile     string `json:"cert_file" env:"CERT_FILE"`
	KeyFile      string `json:"key_file" env:"KEY_FILE"`
	UnixFile     string `json:"unix_file" env:"UNIX_FILE"`
	UnixFilePerm string `json:"unix_file_perm" env:"UNIX_FILE_PERM"`
}

type DataImage struct {
	Dir string `json:"dir" env:"DIR"`
}

type Config struct {
	SiteURL        string `json:"site_url" env:"SITE_URL"`
	JwtSecret      string `json:"jwt_secret" env:"JWT_SECRET"`
	TokenExpiresIn int    `json:"token_expires_in" env:"TOKEN_EXPIRES_IN"`
	DelayedStart   int    `json:"delayed_start" env:"DELAYED_START"`
	Cdn            string `json:"cdn" env:"CDN"`
	DistDir        string `json:"dist_dir"`

	Scheme    Scheme    `json:"scheme"`
	Cors      Cors      `json:"cors" envPrefix:"CORS_"`
	Log       LogConfig `json:"log"`
	Database  Database  `json:"database" envPrefix:"DB_"`
	DataImage DataImage `json:"data_image" envPrefix:"DATA_"`
}

func DefaultConfig() *Config {
	logPath := filepath.Join(flags.DataDir, "log/log.log")
	dbPath := filepath.Join(flags.DataDir, "images.db")
	imagePath := filepath.Join(flags.DataDir, "images")
	// 默认设置
	return &Config{
		JwtSecret:      random.String(16),
		TokenExpiresIn: 48,
		Scheme: Scheme{
			Address:    "0.0.0.0",
			UnixFile:   "",
			HttpPort:   4536,
			HttpsPort:  -1,
			ForceHttps: false,
			CertFile:   "",
			KeyFile:    "",
		},
		Database: Database{
			Type:        "sqlite3",
			Port:        0,
			TablePrefix: "z_",
			DBFile:      dbPath,
		},
		DataImage: DataImage{
			Dir: imagePath,
		},
		Cors: Cors{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		},
		Log: LogConfig{
			Enable:     true,
			Name:       logPath,
			MaxSize:    50,
			MaxBackups: 30,
			MaxAge:     28,
		},
	}
}
