package service

import (
	"com.fs/event-service/event-producer"
	"com.fs/event-service/utils"
	"errors"
	"sync"
)

var producerMapLock sync.Mutex
var producerMap = map[string]event_producer.EventProducer{}

func NewProducer(pluginName string, producerName string, conf string) error {
	if utils.IsStringEmpty(pluginName) || utils.IsStringEmpty(producerName) {
		utils.PrintErr("NewProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := getProducer(producerName)
	if producer != nil {
		utils.PrintErr("NewProducer", "该名称的生产者已存在")
		return errors.New("该名称的生产者已存在")
	}

	p := getProducerPlugin(pluginName)
	if p == nil {
		utils.PrintErr("NewProducer", "没有找到对应的插件")
		return errors.New("没有找到对应的插件")
	}

	producer = p.NewInstance(conf)

	addProducer(producerName, producer)

	return nil
}

func DestroyProducer(pluginName string, producerName string) error {
	if utils.IsStringEmpty(pluginName) || utils.IsStringEmpty(producerName) {
		utils.PrintErr("DestroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := getProducerPlugin(pluginName)
	if p == nil {
		utils.PrintErr("DestroyProducer", "没有找到对应的插件")
		return errors.New("没有找到对应的插件")
	}

	producer := getProducer(producerName)
	if producer == nil {
		utils.PrintErr("DestroyProducer", "没有找到对应的生产者")
		return errors.New("没有找到对应的生产者")
	}

	deleteProducer(producerName)

	p.DestroyInstance(producer)

	return nil
}

func addProducer(producerName string, producer event_producer.EventProducer) {
	producerMapLock.Lock()
	defer producerMapLock.Unlock()

	producerMap[producerName] = producer
}

func getProducer(producerName string) event_producer.EventProducer {
	producerMapLock.Lock()
	defer producerMapLock.Unlock()

	return producerMap[producerName]
}

func deleteProducer(producerName string) {
	producerMapLock.Lock()
	defer producerMapLock.Unlock()

	producerMap[producerName] = nil
}
