package configs

import (
	"os"

	"github.com/joho/godotenv"
)

func InitDotEnv() error {
	var err error
	environment := os.Getenv("ENV")
	if environment == "development" {
		err = godotenv.Load(".env.dev")
	} else {
		err = godotenv.Load()
	}
	return err
}
