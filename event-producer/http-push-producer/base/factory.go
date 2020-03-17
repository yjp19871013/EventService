package base

import (
	"com.fs/event-service/config"
	"com.fs/event-service/event-producer"
	"com.fs/event-service/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
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

type HttpPushFactory struct {
	InitProducer    func(conf *Config) (event_producer.EventProducer, error)
	DestroyProducer func(prod event_producer.EventProducer) error
}

func (factory *HttpPushFactory) NewInstance(producerName string) (event_producer.EventProducer, error) {
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

	pushProducer, err := factory.InitProducer(conf)
	if err != nil {
		utils.PrintCallErr("HttpPushFactory.NewInstance", "producer.InitProducer", err)
		return nil, err
	}

	return pushProducer, nil
}

func (factory *HttpPushFactory) DestroyInstance(prod event_producer.EventProducer) error {
	if prod == nil {
		utils.PrintErr("HttpPushFactory.DestroyInstance", "传递的生产者为nil")
		return errors.New("传递的生产者为nil")
	}

	return factory.DestroyProducer(prod)
}
