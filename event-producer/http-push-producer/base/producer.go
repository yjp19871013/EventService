package base

import (
	"com.fs/event-service/config"
	"com.fs/event-service/utils"
	"context"
	"net/http"
)

type Producer struct {
	Config *Config
	server *http.Server

	OnHandle func(w http.ResponseWriter, r *http.Request)
}

func (prod *Producer) Start() {
	http.HandleFunc(prod.Config.ServerUrl, prod.handleFunc)

	prod.server = &http.Server{
		Addr:    prod.Config.Port,
		Handler: http.DefaultServeMux,
	}

	go func() {
		err := prod.server.ListenAndServe()
		if err != nil {
			utils.PrintCallErr("Producer.Start", "prod.server.ListenAndServe", err)
			return
		}
	}()
}

func (prod *Producer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), config.HttpServerShutdownTimeoutSec)
	defer cancel()

	err := prod.server.Shutdown(ctx)
	if err != nil {
		utils.PrintCallErr("Producer.Stop", "prod.server.Shutdown", err)
		return
	}
}

func (prod *Producer) handleFunc(w http.ResponseWriter, r *http.Request) {
	prod.OnHandle(w, r)
}
