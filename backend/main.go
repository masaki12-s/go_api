package main

import (
	// "github.com/gin-gonic/gin"
	"time"
)

type Task struct {
	ID          uint           `gorm:"primary_key"`
	Task        string         `gorm:"size:255"`
	IsCompleted bool           `gorm:"default:false"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
  }

func main() {

	db := ConnectDB()

	db.AutoMigrate(&Task{})

	// r := gin.Default()

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello World",
	// 	})
	// })
	// r.Run(":8080")
}