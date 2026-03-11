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

	q := url.Values{}
	q.Add("sslmode", os.Getenv("DB_SSLMODE"))

	// Need channel binding param only for neon
	if os.Getenv("APP_ENV") == "production" {
		q.Add("channel_binding", os.Getenv("DB_CHANNEL_BINDING"))
	}

	u := &url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path:     "/" + os.Getenv("DB_NAME"),
		RawQuery: q.Encode(),
	}

	return &Config{
		DBUrl: u.String(),
	}
}
