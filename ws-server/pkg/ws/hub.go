package ws

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client 单个客户端
type Client struct {
	conn      *websocket.Conn
	ChannelID string

	UserID   string
	Username string
}

// Hub 管理中枢
type Hub struct {
	channels map[string]map[*Client]bool
	mu       sync.Mutex
}

type MessagePayload struct {
	ChannelID string    `json:"channel_id"`
	SenderID  string    `json:"sender_id"`
	UserName  string    `json:"username"`
	Message   string    `json:"message"`
	SendTime  time.Time `json:"create_at"`
}

type WSResponse struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewHub() *Hub {
	return &Hub{channels: make(map[string]map[*Client]bool)}
} // 创建空map 用来存频道以及客户端

func (hub *Hub) addClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	// 加个锁 防止冲突

	route := client.ChannelID
	if hub.channels[route] == nil {
		hub.channels[route] = make(map[*Client]bool)
	}
	hub.channels[route][client] = true
}

func (hub *Hub) removeClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	route := client.ChannelID
	delete(hub.channels[route], client)

	if len(hub.channels[route]) == 0 {
		delete(hub.channels, route)
	} // 频道没人了就删整个频道
}

func (hub *Hub) Broadcast(channelID string, response *WSResponse) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for client := range hub.channels[channelID] {
		if jsonBytes, err := json.Marshal(response); err == nil {
			err := client.conn.WriteMessage(websocket.TextMessage, jsonBytes)
			if err != nil {
				return
			}
		} else {
			fmt.Println("marshal error:", err)
		}
	}
}
