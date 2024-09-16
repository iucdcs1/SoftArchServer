package initializers

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		err := godotenv.Load("./.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		logrus.Info("Env variables loaded successfully from .env file!")
	} else {
		logrus.Info("Env variables loaded from GitHub Actions secrets!")
	}
}
