package config

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// http 协议升级

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 可跨域
	},
}
