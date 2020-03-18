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
	"sync"
)

const (
	producerConfigDir = "http_pull"
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
	PullUrl       string `json:"pullUrl"`
	PullPeriodSec uint64 `json:"pullPeriodSec"`
	PullTimeout   uint64 `json:"pullTimeout"`
}

type HttpPullFactory struct {
	sync.Mutex

	InitProducer    func(producerName string, conf *Config) (event_producer.EventProducer, error)
	DestroyProducer func(prod event_producer.EventProducer) error
}

func (factory *HttpPullFactory) NewInstance(producerName string) (event_producer.EventProducer, error) {
	if utils.IsStringEmpty(producerName) {
		utils.PrintErr("HttpPullFactory.NewInstance", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	configFilePath := filepath.Join(config.GetEventServiceConfig().PluginConfig.ProducerConfigDir,
		producerConfigDir, producerName+".json")
	configJson, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		utils.PrintCallErr("HttpPullFactory.NewInstance", "ioutil.ReadFile", err)
		return nil, err
	}

	conf := &Config{}
	err = json.Unmarshal(configJson, conf)
	if err != nil {
		utils.PrintCallErr("HttpPullFactory.NewInstance", "json.Unmarshal", err)
		return nil, err
	}

	pushProducer, err := factory.initProducerWithLock(producerName, conf)
	if err != nil {
		utils.PrintCallErr("HttpPullFactory.NewInstance", "producer.initProducerWithLock", err)
		return nil, err
	}

	return pushProducer, nil
}

func (factory *HttpPullFactory) DestroyInstance(prod event_producer.EventProducer) error {
	if prod == nil {
		utils.PrintErr("HttpPullFactory.DestroyInstance", "传递的生产者为nil")
		return errors.New("传递的生产者为nil")
	}

	return factory.destroyProducerWithLock(prod)
}

func (factory *HttpPullFactory) initProducerWithLock(
	producerName string, conf *Config) (event_producer.EventProducer, error) {
	factory.Lock()
	defer factory.Unlock()

	return factory.InitProducer(producerName, conf)
}

func (factory *HttpPullFactory) destroyProducerWithLock(prod event_producer.EventProducer) error {
	factory.Lock()
	defer factory.Unlock()

	return factory.DestroyProducer(prod)
}
