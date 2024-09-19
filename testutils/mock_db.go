package testutils

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"softarch/initializers"
	"testing"
)

func MockDBSetup(t interface{}) (sqlmock.Sqlmock, *gorm.DB) {
	if t, ok := t.(*testing.T); ok {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %s", err)
		}

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		if err != nil {
			t.Fatalf("failed to open gorm DB: %s", err)
		}

		initializers.DB = gormDB

		return mock, gormDB
	}

	if b, ok := t.(*testing.B); ok {
		db, mock, err := sqlmock.New()
		if err != nil {
			b.Fatalf("failed to open sqlmock database: %s", err)
		}

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		if err != nil {
			b.Fatalf("failed to open gorm DB: %s", err)
		}

		initializers.DB = gormDB

		return mock, gormDB
	}

	panic("Unsupported testing type")
}
