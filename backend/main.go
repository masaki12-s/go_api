package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Stock struct {
	ID		uint	`json:"id"`
	Name	string	`json:"name"`
	Price	float64	`json:"price"`
	Amount	int		`json:"amount"`
  }

func main() {

	db := ConnectDB()
	db.AutoMigrate(&Stock{})

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	// 在庫の更新, 作成
    /**
     * POST /v1/stocks
     * http resuest body
     * {
     * "name": "string" (required),
     * "amount": 0 (optional),
     * }
     * 
     * http response header
     * Location: /v1/stock/:name
     * 
     * http response body
     * same as request body
    **/
	r.POST("/v1/stocks", func(c *gin.Context) {
		var stock Stock
		if err := c.ShouldBindJSON(&stock); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var exsitingStock Stock
		result := db.Where("name = ?", stock.Name).First(&exsitingStock)
		
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				if err := db.Create(&stock).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
					return
				}
			}
		} else {
			new_amount := exsitingStock.Amount + stock.Amount
			db.Model(&exsitingStock).Update("amount", new_amount)
		}
		c.JSON(http.StatusOK, exsitingStock)
	})

	r.Run(":8080")
}
