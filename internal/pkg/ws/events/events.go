package events

import (
	"context"
	"encoding/json"
	"fmt"
	rediscli "ws_server/internal/pkg/redis"
	"ws_server/internal/pkg/ws"
)

type Envelope struct {
	TaskType string          `json:"task_type"`
	Payload  json.RawMessage `json:"payload"`
}

// 订阅主服务发送的频道加入事件
func ListenChannelEvents(hub *ws.Hub) {
	ctx := context.Background()
	sub := rediscli.Rds.Subscribe(ctx, "channel_event")
	ch := sub.Channel()

	for msg := range ch {
		var env Envelope
		if err := json.Unmarshal([]byte(msg.Payload), &env); err != nil {
			fmt.Println("invalid envelope:", err)
			continue
		}

		switch env.TaskType {
		case "JOIN":
			handleJoinEvent(hub, env.Payload)
		case "LEAVE":
			handleLeaveEvent(hub, env.Payload)
		case "KICK":
			handleKickEvent(hub, env.Payload)
		case "CHANGE":
			handleChangeEvent(hub, env.Payload)
		case "DELETE":
			handleDeleteEvent(hub, env.Payload)
		}
	}
}
