package main

import (
	"encoding/json"
	"fmt"
	"go-order-event-service/internal/event"
	"go-order-event-service/internal/stream"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ✅ 建立全域 Hub（只建一次）
	hub := stream.NewHub()

	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// ===== SSE 連線 =====
	r.GET("/stream/orders", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "user_id required"})
			return
		}

		// ⭐⭐⭐ CORS（SSE 一定要在這）
		// ⭐⭐⭐ CORS（SSE 一定要在這）
		allowedOrigin := os.Getenv("SSE_ALLOW_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:5173" // 預設給本機
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(500, gin.H{"error": "streaming unsupported"})
			return
		}

		client := &stream.Client{
			UserID: userID,
			Ch:     make(chan []byte, 10),
		}

		hub.Add(userID, client)
		defer hub.Remove(userID, client)

		notify := c.Writer.CloseNotify()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg := <-client.Ch:
				fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
				flusher.Flush()

			case <-ticker.C:
				// ⭐ SSE 心跳（不影響前端）
				fmt.Fprintf(c.Writer, "event: ping\ndata: {}\n\n")
				flusher.Flush()

			case <-notify:
				return
			}
		}

	})

	// ===== 接收 Laravel 事件 =====
	r.POST("/events/order", func(c *gin.Context) {
		var env event.Envelope
		if err := c.ShouldBindJSON(&env); err != nil {
			c.JSON(400, gin.H{"error": "invalid event format"})
			return
		}

		log.Printf("[EVENT] %s from %s", env.Event, env.Source)

		// ⭐ 先只處理 order.created
		if env.Event == "order.created" {
			payload, _ := json.Marshal(env)

			dataMap, ok := env.Data.(map[string]interface{})
			if ok {
				userID := fmt.Sprint(dataMap["user_id"])
				hub.Publish(userID, payload)
			}
		}

		c.JSON(200, gin.H{"ok": true})
	})

	r.Run(":8080")
}
