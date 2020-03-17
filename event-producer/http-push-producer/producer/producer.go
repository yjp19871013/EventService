package producer

import (
	"com.fs/event-service/config"
	"com.fs/event-service/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	producerConfigDir = "http_push"
)

func init() {
	conf := config.GetEventServiceConfig().PluginConfig
	configDir := filepath.Join(conf.ProducerConfigDir, producerConfigDir)
	if !utils.PathExists(configDir) {
		err := os.MkdirAll(configDir, os.ModePerm|os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
}

type Config struct {
	ServerUrl string `json:"serverUrl"`
	Port      string `json:"port"`
	Method    string `json:"method"`
}

type Producer struct {
	config *Config
	server *http.Server

	OnHandle func(w http.ResponseWriter, r *http.Request)
}

func InitProducer(producerName string) (*Producer, error) {
	if utils.IsStringEmpty(producerName) {
		utils.PrintErr("InitProducer", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	configFilePath := filepath.Join(config.GetEventServiceConfig().PluginConfig.ProducerConfigDir,
		producerConfigDir, producerName+".json")
	configJson, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		utils.PrintCallErr("InitProducer", "ioutil.ReadFile", err)
		return nil, err
	}

	conf := &Config{}
	err = json.Unmarshal(configJson, conf)
	if err != nil {
		utils.PrintCallErr("InitProducer", "json.Unmarshal", err)
		return nil, err
	}

	prod := &Producer{config: conf}

	return prod, nil
}

func DestroyProducer(prod *Producer) {
	if prod == nil {
		return
	}

	prod.config = nil
	prod = nil
}

func (prod *Producer) Start() {
	http.HandleFunc(prod.config.ServerUrl, prod.handleFunc)

	prod.server = &http.Server{
		Addr:    prod.config.Port,
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
	fmt.Println("handleFunc Here")
	prod.OnHandle(w, r)
}
