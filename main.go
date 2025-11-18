package main

import (
	"fmt"
	"net/http"

	"ws_server/internal/pkg/ws"
	"ws_server/internal/router"
)

func main() {
	hub := ws.NewHub()

	r := router.SetupWSRouter(hub)

	fmt.Println("try open port")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
