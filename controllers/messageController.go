package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"softarch/initializers"
	"softarch/models"
	"time"
)

type messageData struct {
	Date time.Time `json:"created_at"`
	Text string    `json:"text"`
}

func (d messageData) MarshalJSON() ([]byte, error) {
	type Alias messageData
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result []messageData
	for _, msg := range messages {
		result = append(result, messageData{Date: msg.Date, Text: msg.Text})
	}

	c.IndentedJSON(http.StatusOK, result)
}

func SendMessage(c *gin.Context) {
	var message messageData

	if err := c.BindJSON(&message); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newMessage := models.Message{
		Text: message.Text,
		Date: time.Now(),
	}

	result := initializers.DB.Model(&models.Message{}).Create(&newMessage)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to send message"})
		return
	}

	responseMessage := messageData{
		Date: newMessage.Date,
		Text: newMessage.Text,
	}

	c.IndentedJSON(http.StatusOK, responseMessage)
}

func GetCount(c *gin.Context) {
	var messages []messageData
	err := initializers.DB.Model(&models.Message{}).Find(&messages).Error

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, len(messages))
}
