package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/zechao158/ecomm/storage"
)

type Config struct {
	APPEnv   string
	HTTPHost string
	HTTPPort string
	storage.Config
}

var ENVs = initConfig()

func initConfig() Config {
	appEnv := "local"
	if e := os.Getenv("APP_ENV"); e != "" {
		appEnv = e
	}
	if appEnv == "local" {
		godotenv.Load()
	}

	debug, _ := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	return Config{
		APPEnv:   appEnv,
		HTTPHost: getEnv("HTTP_HOST", "localhost"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		Config: storage.Config{
			DBUser:     getEnv("DB_USER", "ecom"),
			DBName:     getEnv("DB_NAME", "ecom"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPassword: getEnv("DB_PASSWORD", "ecom"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
			DebugMode:  debug,
		},
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}