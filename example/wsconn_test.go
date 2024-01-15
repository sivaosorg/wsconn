package example

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sivaosorg/govm/wsconnx"
	"github.com/sivaosorg/wsconn"
)

func TestWebsocket(t *testing.T) {
	r := gin.Default()
	ws := wsconn.NewWebsocket()
	s := wsconn.NewWebsocketService(ws)

	r.GET("/subscribe", s.SubscribeMessage) // ws://localhost:8081/subscribe
	r.POST("/message", func(c *gin.Context) {
		var message wsconnx.WsConnMessagePayload
		message.SetGenesisTimestamp(time.Now())
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		s.BroadcastMessage(message)
		c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully", "data": message})
	})
	r.Run(":8081")
}
