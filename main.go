package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// 主服务
	"simplechat/server/common/async"
	"simplechat/server/common/pkg/db"
	rediscli "simplechat/server/common/pkg/redis"
	serverRouter "simplechat/server/router"

	// WS 服务

	wsPkg "simplechat/ws-server/pkg/ws"
	wsEvents "simplechat/ws-server/pkg/ws/events"
	wsRouter "simplechat/ws-server/router"
)

func main() {
	_ = godotenv.Load(".env")

	db.Init()
	rediscli.Init()

	// hub
	hub := wsPkg.NewHub()
	go wsEvents.ListenChannelEvents(hub)

	// 异步 worker
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go async.StartWorker(ctx)

	// 全局 Gin
	r := gin.Default()

	// 注册 API 路由
	serverRouter.SetupRouter(r)

	// 注册 WS 路由
	wsRouter.SetupWSRouter(r, hub)

	// 统一端口
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Unified Server started on :8080")

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
