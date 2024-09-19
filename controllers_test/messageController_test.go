package controllers_test

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"softarch/controllers"
	"softarch/initializers"
	"strings"
	"testing"
	"time"
)

func mockDBSetup(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
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

func TestGetMessages(t *testing.T) {
	mock, _ := mockDBSetup(t)

	rows := sqlmock.NewRows([]string{"created_at", "text"}).
		AddRow(time.Now(), "Test message 1").
		AddRow(time.Now(), "Test message 2")

	mock.ExpectQuery("SELECT \\* FROM \"messages\"").
		WillReturnRows(rows)

	router := gin.Default()
	router.GET("/messages", controllers.GetMessages)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []controllers.ResultingMessageData
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Test message 1", result[0].Text)
	assert.Equal(t, "Test message 2", result[1].Text)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSendMessage(t *testing.T) {
	mock, _ := mockDBSetup(t)

	var expectedSQL = `INSERT INTO "messages" \("date","text"\) VALUES \(\$1,\$2\) RETURNING "id"`

	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQL).
		WithArgs(sqlmock.AnyArg(), "Test message").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	router := gin.Default()
	router.POST("/messages", controllers.SendMessage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", strings.NewReader(`{"text": "Test message"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.ResultingMessageData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test message", response.Text)
	assert.NotEmpty(t, response.Date)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetCount(t *testing.T) {
	mock, _ := mockDBSetup(t)

	mock.ExpectQuery(`SELECT count\(\*\) FROM "messages"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	router := gin.Default()
	router.GET("/messages/count", controllers.GetCount)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages/count", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, "2", w.Body.String())

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
