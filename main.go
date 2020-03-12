package main

import (
	"com.fs/event-service/config"
	"com.fs/event-service/router"
	"context"
	"github.com/gin-gonic/gin"
	DEATH "gopkg.in/vrecan/death.v3"
	"log"
	"net/http"
	"syscall"
	"time"
)

// @title 事件回调微服务
// @version 0.1
// @BasePath /
func main() {
	r := gin.Default()

	router.InitRouter(r)

	srv := &http.Server{
		Addr:    config.GetEventServiceConfig().ServerConfig.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	death := DEATH.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	_ = death.WaitForDeath()
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
