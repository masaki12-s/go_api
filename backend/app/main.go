package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Stock struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Amount int     `json:"amount"`
}

type checkedStock struct {
	Name   string  `json:"name"`
	Amount int     `json:"amount"`
}

func updateStock(c *gin.Context, db *gorm.DB) {
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
}

func checkStock(c *gin.Context, db *gorm.DB) {
	var stocks []Stock
	var checkedstocks []checkedStock
	var result *gorm.DB
	if name := c.Param("name"); name == "" {
		result = db.Find(&stocks).Find(&checkedstocks)
	} else {
		result = db.Where("name = ?", name).Find(&stocks).Find(&checkedstocks)
	}
	// エラー処理
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error!"})
		return
	}
	// 正常時処理
	c.JSON(http.StatusOK, checkedstocks)
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
        updateStock(c, db)
	})

    // 在庫チェック
    /**
    * GET /v1/stocks(/:name)
    * http response body
    * :nameが指定されている場合
    * {[name]: [amount]} (在庫がない場合はamountは0)
    * :nameが指定されていない場合
    * 全ての商品の在庫の数を:nameで昇順sortして
    * {[name]: [amount], [name]: [amount], ...}で返す
    * 例: {"apple": 1, "banana": 2, "orange": 0}
    * 在庫が0のものは含まない
    **/
    r.GET("/v1/stocks", func(c *gin.Context) {
        checkStock(c, db)
    })
	r.GET("/v1/stocks/:name", func(c *gin.Context) {
		checkStock(c, db)
	})

	r.Run(":8080")
}



// 販売
/**
 * POST /v1/sales
 * http request body
 * {
 * "name": "string" (required),
 * "amount": 0 (optional), 対象の商品を在庫から販売する数（正の整数）を指定する。省略時は 1 とする。
 * "price": 0 (optional), 対象の商品の価格（0より大きい数値）を指定する。入力された時のみ、price x amount を売り上げに加算する。
 * }
 *
 * http response header
 * Location: /v1/sales/:name
 *
 * http response body
 * request bodyと同じ
 **/

// 売り上げチェック
/**
 * GET /v1/sales
 *
 * http response body
 * {
 * "sales": 0 (売り上げの合計額, 小数第二位まで, 整数の場合は小数第一位まで)
 * }
 *
 **/

// 全削除
/**
 * DELETE /v1/stocks
**/
