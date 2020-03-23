package base

import (
	"com.fs/event-service/config"
	"com.fs/event-service/utils"
	"context"
	"errors"
	"net/http"
	"strings"
)

type Producer struct {
	Config *Config

	OnHandle func(w http.ResponseWriter, r *http.Request)

	server *http.Server
}

func (prod *Producer) Start() error {
	if prod.Config == nil {
		utils.PrintErr("Producer.Start", "没有传递配置")
		return errors.New("没有传递配置")
	}

	if prod.OnHandle == nil {
		utils.PrintErr("Producer.Start", "没有创建OnHandle函数")
		return errors.New("没有创建OnHandle函数")
	}

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

	return nil
}

func (prod *Producer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.HttpServerShutdownTimeoutSec)
	defer cancel()

	err := prod.server.Shutdown(ctx)
	if err != nil {
		utils.PrintCallErr("Producer.Stop", "prod.server.Shutdown", err)
		return err
	}

	return nil
}

func (prod *Producer) handleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != strings.ToUpper(prod.Config.Method) {
		return
	}

	prod.OnHandle(w, r)
}
