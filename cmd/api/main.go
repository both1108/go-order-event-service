package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 健康檢查（部署一定會用到）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 未來 Laravel 會打這個
	r.POST("/events/order", func(c *gin.Context) {
		var payload map[string]any
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// 先只 log，後面再處理事件
		c.JSON(http.StatusOK, gin.H{
			"received": true,
			"payload":  payload,
		})
	})

	r.Run(":8080") // 本機 port
}
