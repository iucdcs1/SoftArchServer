package initializers

import (
	"github.com/sirupsen/logrus"
	"log"
	"softarch/models"
)

func SyncDatabase() {
	migrateModel(&models.Message{})

	logrus.Info("Database migration completed!")
}

func migrateModel(model interface{}) {
	if !DB.Migrator().HasTable(model) {
		if err := DB.AutoMigrate(model); err != nil {
			log.Fatalf("Database migration for %T failed: %v", model, err)
		}
	} else {
		logrus.Infof("Skipping migration for %T - table already exists", model)
	}
}
