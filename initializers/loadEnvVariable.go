package initializers

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		err := godotenv.Load("./.env")
		if err != nil {
			err := godotenv.Load("../.env")
			if err != nil {
				logrus.Warn("Error loading .env file: ", err)
			}
		} else {
			logrus.Info("Env variables loaded successfully from .env file!")
		}
	} else {
		logrus.Info("Env variables loaded from GitHub Actions secrets!")
	}
}
