package edge_cases_tests

import (
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"softarch/controllers"
	"softarch/testutils"
	"strings"
	"testing"
)

func TestGetMessagesNoMessages(t *testing.T) {
	mock, _ := testutils.MockDBSetup(t)

	mock.ExpectQuery("SELECT \\* FROM \"messages\"").
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "text"}))

	router := gin.Default()
	router.GET("/messages", controllers.GetMessages)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []controllers.ResultingMessageData
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Empty(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSendMessageInvalidJSON(t *testing.T) {
	mock, _ := testutils.MockDBSetup(t)

	router := gin.Default()
	router.POST("/messages", controllers.SendMessage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", strings.NewReader(`{"text": `)) // Invalid JSON
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid input")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetMessagesDBError(t *testing.T) {
	mock, _ := testutils.MockDBSetup(t)

	mock.ExpectQuery("SELECT \\* FROM \"messages\"").
		WillReturnError(errors.New("database error"))

	router := gin.Default()
	router.GET("/messages", controllers.GetMessages)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "database error")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSendMessageIncorrectDataType(t *testing.T) {
	mock, _ := testutils.MockDBSetup(t)

	router := gin.Default()
	router.POST("/messages", controllers.SendMessage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", strings.NewReader(`{"text": 12345}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid input")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSendMessageLargePayload(t *testing.T) {
	mock, _ := testutils.MockDBSetup(t)

	longText := strings.Repeat("A", 10000) // Large payload
	var expectedSQL = `INSERT INTO "messages" \("date","text"\) VALUES \(\$1,\$2\) RETURNING "id"`

	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQL).
		WithArgs(sqlmock.AnyArg(), longText).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	router := gin.Default()
	router.POST("/messages", controllers.SendMessage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", strings.NewReader(`{"text": "`+longText+`"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.ResultingMessageData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, longText, response.Text)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSendMessageDuplicate(t *testing.T) {
	mock, _ := testutils.MockDBSetup(t)

	messageText := "Duplicate message"
	var expectedSQL = `INSERT INTO "messages" \("date","text"\) VALUES \(\$1,\$2\) RETURNING "id"`

	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQL).
		WithArgs(sqlmock.AnyArg(), messageText).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	router := gin.Default()
	router.POST("/messages", controllers.SendMessage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", strings.NewReader(`{"text": "`+messageText+`"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.ResultingMessageData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, messageText, response.Text)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
