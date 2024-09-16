package models

import (
	"time"
)

type Message struct {
	ID   uint      `gorm:"primary_key;auto_increment" json:"id"`
	Date time.Time `sql:"type:timestamp without time zone" json:"created_at"`
	Text string    `gorm:"not null" json:"text"`
}
