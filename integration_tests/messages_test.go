package integration_tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"softarch/controllers"
	"softarch/initializers"
	"softarch/models"
	"testing"
	"time"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/messages", controllers.GetMessages)
	router.POST("/messages", controllers.SendMessage)
	return router
}

func setupDatabase() *gorm.DB {
	initializers.LoadEnvVariables()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: os.Getenv("DATABASE_URL"),
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func TestGetMessagesIntegration(t *testing.T) {
	router := setupRouter()
	db := setupDatabase()
	initializers.DB = db

	var initialMessageCount int64
	db.Model(&models.Message{}).Count(&initialMessageCount)

	db.Create(&models.Message{Text: "Test1_test", Date: time.Now()})
	db.Create(&models.Message{Text: "Test2_test", Date: time.Now()})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var messages []controllers.ResultingMessageData
	err := json.Unmarshal(w.Body.Bytes(), &messages)
	assert.NoError(t, err)

	assert.Equal(t, int(initialMessageCount+2), len(messages))
}

func TestSendMessageIntegration(t *testing.T) {
	router := setupRouter()
	db := setupDatabase()
	initializers.DB = db

	var initialMessageCount int64
	db.Model(&models.Message{}).Count(&initialMessageCount)

	message := controllers.MessageData{
		Text: "Test message",
	}

	messageModelled := models.Message{
		Date: time.Now(),
		Text: message.Text,
	}

	payload, _ := json.Marshal(messageModelled)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.ResultingMessageData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, message.Text, response.Text)

	var dbMessage models.Message
	db.First(&dbMessage, "text = ?", message.Text)
	assert.Equal(t, message.Text, dbMessage.Text)

	var newMessageCount int64
	db.Model(&models.Message{}).Count(&newMessageCount)
	assert.Equal(t, initialMessageCount+1, newMessageCount)
}
