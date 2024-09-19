package initializers_test

import (
	"os"
	"softarch/initializers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvVariables(t *testing.T) {
	err := os.Setenv("DATABASE_URL", "test-db-url")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer func() {
		err := os.Unsetenv("DATABASE_URL")
		if err != nil {
			assert.Error(t, err)
		}
	}()

	initializers.LoadEnvVariables()

	assert.Equal(t, "test-db-url", os.Getenv("DATABASE_URL"))
}
