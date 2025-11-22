package main

import (
	"fmt"
	"net/http"

	rediscli "ws_server/internal/pkg/redis"
	"ws_server/internal/pkg/ws"
	"ws_server/internal/pkg/ws/events"
	"ws_server/internal/router"
)

func main() {
	hub := ws.NewHub()

	rediscli.Init()
	go events.ListenChannelEvents(hub)

	r := router.SetupWSRouter(hub)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
