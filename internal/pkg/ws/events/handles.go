package events

import (
	"encoding/json"
	"fmt"
	"ws_server/internal/pkg/ws"
)

type JoinChannelPayload struct {
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	IsOwner     bool   `json:"is_owner"`
	UserID      string `json:"user_id"`
}

type LeaveChannelPayload struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}

type ChangeChannelNamePayload struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

type DeleteChannelPayload struct {
	ChannelID string `json:"channel_id"`
}

func handleJoinEvent(hub *ws.Hub, raw json.RawMessage) {
	var p JoinChannelPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		fmt.Println("JOIN event unmarshal failed:", err)
		return
	}

	hub.Broadcast(p.ChannelID, &ws.WSResponse{
		Type:    "JOIN",
		Payload: p,
	})
}

func handleLeaveEvent(hub *ws.Hub, raw json.RawMessage) {
	var p LeaveChannelPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		fmt.Println("LEAVE event unmarshal failed:", err)
		return
	}

	hub.Broadcast(p.ChannelID, &ws.WSResponse{
		Type:    "LEAVE",
		Payload: p,
	})
}

func handleKickEvent(hub *ws.Hub, raw json.RawMessage) {
	var p LeaveChannelPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		fmt.Println("Kick event unmarshal failed:", err)
		return
	}

	hub.Broadcast(p.ChannelID, &ws.WSResponse{
		Type:    "KICK",
		Payload: p,
	})
}

func handleChangeEvent(hub *ws.Hub, raw json.RawMessage) {
	var p ChangeChannelNamePayload
	if err := json.Unmarshal(raw, &p); err != nil {
		fmt.Println("Change event unmarshal failed:", err)
		return
	}

	hub.Broadcast(p.ChannelID, &ws.WSResponse{
		Type:    "CHANGE",
		Payload: p,
	})
}

func handleDeleteEvent(hub *ws.Hub, raw json.RawMessage) {
	var p DeleteChannelPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		fmt.Println("Delete event unmarshal failed:", err)
		return
	}

	hub.Broadcast(p.ChannelID, &ws.WSResponse{
		Type:    "DELETE",
		Payload: p,
	})
}
