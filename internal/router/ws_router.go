package router

import (
	"net/http"

	"ws_server/internal/pkg/ws"

	"github.com/gorilla/mux"
)

// 这里用的 gorilla ws 配套的 mux, 喜欢gin的话 可以改用过去, 最好用gin成一套体系
func SetupWSRouter(hub *ws.Hub) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ws/chat", func(w http.ResponseWriter, r *http.Request) {
		ws.HandlerWebSocket(hub, w, r)
	}).Methods("GET") // ws chat 路由, 参数 channelID, username, token

	router.HandleFunc("/ws/aibot", func(w http.ResponseWriter, r *http.Request) {
		ws.HandlerWebSocket(hub, w, r)
	}).Methods("GET") // 看看吧 用ai拓展那边使用?

	return router
}
