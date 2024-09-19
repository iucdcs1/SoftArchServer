package initializers

import (
	"github.com/sirupsen/logrus"
	"softarch/models"
)

func SyncDatabase() {
	if err := DB.AutoMigrate(&models.Message{}); err != nil {
		logrus.Fatalf("Error migrating database: %v", err)
	}

	logrus.Info("Database migration completed!")
}
