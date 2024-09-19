package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"softarch/initializers"
	"softarch/models"
	"time"
)

type MessageData struct {
	Date time.Time `json:"created_at"`
	Text string    `json:"text"`
}

type ResultingMessageData struct {
	Date string `json:"created_at"`
	Text string `json:"text"`
}

func (d MessageData) MarshalJSON() ([]byte, error) {
	type Alias MessageData
	return json.Marshal(&struct {
		Date string `json:"created_at"`
		*Alias
	}{
		Date:  d.Date.Format("02.01.2006 15:04"),
		Alias: (*Alias)(&d),
	})
}

func GetMessages(c *gin.Context) {
	var messages []models.Message

	err := initializers.DB.Model(&models.Message{}).Find(&messages).Error

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []ResultingMessageData
	for _, msg := range messages {
		result = append(result, ResultingMessageData{
			Date: msg.Date.Format("02.01.2006 15:04"),
			Text: msg.Text,
		})
	}

	c.IndentedJSON(http.StatusOK, result)
}

func SendMessage(c *gin.Context) {
	var message MessageData

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if message.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text field is required"})
		return
	}

	newMessage := models.Message{
		Text: message.Text,
		Date: time.Now(),
	}

	result := initializers.DB.Model(&models.Message{}).Create(&newMessage)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to send message: " + result.Error.Error()})
		return
	}

	responseMessage := MessageData{
		Date: newMessage.Date,
		Text: newMessage.Text,
	}

	c.IndentedJSON(http.StatusOK, responseMessage)
}

func GetCount(c *gin.Context) {
	var count int64
	err := initializers.DB.Model(&models.Message{}).Count(&count).Error

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, count)
}
