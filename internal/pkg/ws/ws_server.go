package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"ws_server/internal/config"
)

func HandlerWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {

	channelId := r.URL.Query().Get("channelID") // 获取此客户端的所在频道, 用于发送信息
	username := r.URL.Query().Get("username")   // 获取用户名当sender
	//token := r.URL.Query().Get("token")         // 获取用户名当sender

	//// token 校验
	//if !jwt.exists(token) {
	//	fmt.Println("token not exists")
	//	http.Error(w, "access token does not exist", http.StatusUnauthorized)
	//	return
	//}
	//
	//// 数据库校验
	//if !db.ChannelExists(channelId) { // 数据库找channel ID, 没找到不能进行 ws 连接
	//	fmt.Println("db cannot serch this channel ID:", err)
	//	http.Error(w, "channel not found", http.StatusForbidden)
	//	return
	//}

	// 升级协议, 并支持跨域
	conn, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		http.Error(w, "upgrade error", http.StatusInternalServerError)
		return
	}

	client := &Client{
		conn:      conn,
		channelId: channelId,
		username:  username,
	}

	hub.addClient(client) // 加入客户端到 hub
	defer func() {
		hub.removeClient(client)
		conn.Close()
		fmt.Println("Client disconnected:", client.username)
	}()

	fmt.Println(client.username, "channel:", client.channelId, "-----------connected")

	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil { // 连接断开 直接退出循环
			return
		}

		if msgType != websocket.TextMessage { // 如果不是文本直接下一个循环, 之处理文本
			continue
		}

		ResponseMessage := &WSMessageResponse{
			Type: "SYSTEM",
			Data: WSMessage{
				ChannelId: client.channelId,
				Username:  client.username,
				Message:   "server received: " + string(message),
				SendTime:  time.Now().UTC(),
			},
		}

		// 服务器回复请求成功
		if jsonBytes, err := json.Marshal(ResponseMessage); err == nil {
			err := conn.WriteMessage(websocket.TextMessage, jsonBytes)
			if err != nil {
				return
			}
		} else {
			fmt.Println("marshal error:", err)
		}

		ResponseMessage.Type = "CHAT"
		ResponseMessage.Data.Message = string(message)
		hub.broadcast(client.channelId, ResponseMessage) // 广播消息
	}
}
