package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"simplechat/server/common/pkg/jwt"
	rediscli "simplechat/server/common/pkg/redis"
	"strings"
	"time"

	"simplechat/ws-server/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const StreamName = "async_tasks"

type Envelope struct {
	TaskType string          `json:"task_type"`
	Payload  json.RawMessage `json:"payload"`
}

type Message struct {
	SenderID string    `json:"sender_id"`
	Sender   string    `json:"sender"`
	Content  string    `json:"content"`
	SentAt   time.Time `json:"sent_at"`
}

func HandlerWebSocketGin(hub *Hub, c *gin.Context) {
	// 从 gin.Context 拿出 ResponseWriter 和 Request
	w := c.Writer
	r := c.Request

	// 调用你原来的处理函数
	HandlerWebSocket(hub, w, r)
}

func HandlerWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {

	ChannelID := r.URL.Query().Get("channel_id") // 获取此客户端的所在频道, 用于发送信息
	UserID := r.URL.Query().Get("user_id")       // 获取用户名当sender
	Username := r.URL.Query().Get("username")    // 获取用户名当sender
	token := r.URL.Query().Get("token")          // 获取用户名当sender

	// token 校验
	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Printf("Invalid Authorization format")
		return
	}

	_, err := jwt.ParseToken(parts[1])
	if err != nil {
		log.Printf("Invalid or expired token")
		return
	}

	// 升级协议, 并支持跨域
	conn, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		http.Error(w, "upgrade error", http.StatusInternalServerError)
		return
	}

	client := &Client{
		conn:      conn,
		ChannelID: ChannelID,
		UserID:    UserID,
		Username:  Username,
	}

	hub.addClient(client) // 加入客户端到 hub
	defer func() {
		hub.removeClient(client)
		conn.Close()
		fmt.Println("Client disconnected:", client.UserID)
	}()

	fmt.Println(client.UserID, "channel:", client.ChannelID, "-----------connected")

	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil { // 连接断开 直接退出循环
			return
		}

		if msgType != websocket.TextMessage { // 如果不是文本直接下一个循环, 之处理文本
			continue
		}

		ResponseMessage := &WSResponse{
			Type: "SYSTEM",
			Payload: MessagePayload{
				ChannelID: client.ChannelID,
				SenderID:  client.UserID,
				UserName:  Username,
				Message:   string(message),
				SendTime:  time.Now().UTC(),
			},
		}

		// 服务器回复请求成功
		jsonBytes, err := json.Marshal(ResponseMessage)
		if err == nil {
			err := conn.WriteMessage(websocket.TextMessage, jsonBytes)
			if err != nil {
				return
			}
		} else {
			fmt.Println("marshal error:", err)
		}

		mp, _ := ResponseMessage.Payload.(MessagePayload)
		msg := &Message{
			SenderID: client.UserID,
			Content:  string(message),
			SentAt:   mp.SendTime,
		}

		msgJSON, _ := json.Marshal(msg)
		rediscli.Rds.LPush(context.Background(), fmt.Sprintf("chat:history:%s", client.ChannelID), msgJSON)
		rediscli.Rds.LTrim(context.Background(), fmt.Sprintf("chat:history:%s", client.ChannelID), 0, 99)

		ResponseMessage.Type = "CHAT"
		payloadBytes, err := json.Marshal(ResponseMessage.Payload)
		if err != nil {
			return
		}

		envelope := Envelope{
			TaskType: "chat_message",
			Payload:  payloadBytes,
		}
		envelopeBytes, err := json.Marshal(envelope)
		if err != nil {
			return
		}

		_, err = rediscli.Rds.XAdd(context.Background(), &redis.XAddArgs{
			Stream: StreamName,
			Values: map[string]interface{}{
				"data": string(envelopeBytes),
			},
		}).Result()

		if err != nil {
			log.Printf("Failed to enqueue task", err)
			return
		}

		hub.Broadcast(client.ChannelID, ResponseMessage) // 广播消息
	}
}
