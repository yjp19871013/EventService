package base

import (
	"com.fs/event-service/config"
	"com.fs/event-service/plugins"
	"com.fs/event-service/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"sync"
)

const (
	httpPushInstanceDir = "http_push"
)

type Config struct {
	ServerUrl string `json:"serverUrl"`
	Port      string `json:"port"`
	Method    string `json:"method"`
}

type HttpPushFactory struct {
	sync.Mutex

	InitProducer         func(conf *Config) (plugins.Instance, error)
	DestroyProducer      func(instance plugins.Instance) error
	OfferInstancesSubDir func() string
}

func (factory *HttpPushFactory) NewInstance(instanceName string) (plugins.Instance, error) {
	if utils.IsStringEmpty(instanceName) {
		utils.PrintErr("HttpPushFactory.NewInstance", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	if factory.InitProducer == nil {
		utils.PrintErr("HttpPushFactory.NewInstance", "没有传递InitProducer")
		return nil, errors.New("没有传递InitProducer")
	}

	if factory.DestroyProducer == nil {
		utils.PrintErr("HttpPushFactory.NewInstance", "没有传递DestroyProducer")
		return nil, errors.New("没有传递DestroyProducer")
	}

	if factory.OfferInstancesSubDir == nil {
		utils.PrintErr("HttpPushFactory.NewInstance", "没有传递OfferInstancesSubDir")
		return nil, errors.New("没有传递OfferInstancesSubDir")
	}

	instanceFilePath := filepath.Join(factory.GetInstancesDir(), instanceName+".json")
	configJson, err := ioutil.ReadFile(instanceFilePath)
	if err != nil {
		utils.PrintCallErr("HttpPushFactory.NewInstance", "ioutil.ReadFile", err)
		return nil, err
	}

	conf := &Config{}
	err = json.Unmarshal(configJson, conf)
	if err != nil {
		utils.PrintCallErr("HttpPushFactory.NewInstance", "json.Unmarshal", err)
		return nil, err
	}

	pushProducer, err := factory.initInstanceWithLock(conf)
	if err != nil {
		utils.PrintCallErr("HttpPushFactory.NewInstance", "producer.initProducerWithLock", err)
		return nil, err
	}

	return pushProducer, nil
}

func (factory *HttpPushFactory) DestroyInstance(instance plugins.Instance) error {
	if instance == nil {
		utils.PrintErr("HttpPushFactory.DestroyInstance", "传递的生产者为nil")
		return errors.New("传递的生产者为nil")
	}

	return factory.destroyProducerWithLock(instance)
}

func (factory *HttpPushFactory) GetInstancesDir() string {
	conf := config.GetEventServiceConfig().PluginConfig
	return filepath.Join(conf.ProducerConfigDir, httpPushInstanceDir, factory.OfferInstancesSubDir())
}

func (factory *HttpPushFactory) initInstanceWithLock(conf *Config) (plugins.Instance, error) {
	factory.Lock()
	defer factory.Unlock()

	return factory.InitProducer(conf)
}

func (factory *HttpPushFactory) destroyProducerWithLock(instance plugins.Instance) error {
	factory.Lock()
	defer factory.Unlock()

	return factory.DestroyProducer(instance)
}
