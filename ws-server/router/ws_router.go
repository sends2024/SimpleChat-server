package router

import (
	"simplechat/ws-server/pkg/ws"

	"github.com/gin-gonic/gin"
)

func SetupWSRouter(r *gin.Engine, hub *ws.Hub) {
	r.GET("/ws/chat", func(c *gin.Context) {
		ws.HandlerWebSocketGin(hub, c)
	})
}
