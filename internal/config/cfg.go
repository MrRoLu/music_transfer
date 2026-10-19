package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DevToken   string
	MediaUserToken string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}

	return &Config{
		DevToken:   os.Getenv("DEV_TOKEN_APPLE"),
		MediaUserToken: os.Getenv("MEDIA_USER_TOKEN_APPLE"),
	}
}
