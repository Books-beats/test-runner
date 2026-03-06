package config

import (
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path:     "/" + os.Getenv("DB_NAME"),
		RawQuery: "sslmode=" + os.Getenv("DB_SSLMODE") + "&pgbouncer=true&pool_timeout=10",
	}

	return &Config{
		DBUrl: u.String(),
	}
}
