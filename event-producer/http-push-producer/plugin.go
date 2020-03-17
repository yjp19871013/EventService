package main

import (
	"com.fs/event-service/event-producer"
	"com.fs/event-service/event-producer/http-push-producer/producer"
	"com.fs/event-service/utils"
	"errors"
	"fmt"
)

func init() {
	fmt.Println("http push producer load")
}

type HttpPushPlugin struct {
}

var Plugin = HttpPushPlugin{}

func (p *HttpPushPlugin) NewInstance(conf string) (event_producer.EventProducer, error) {
	pushProducer, err := producer.InitProducer(conf)
	if err != nil {
		utils.PrintCallErr("HttpPushPlugin.NewInstance", "producer.InitProducer", err)
		return nil, err
	}

	return pushProducer, nil
}

func (p *HttpPushPlugin) DestroyInstance(prod event_producer.EventProducer) error {
	if prod == nil {
		utils.PrintErr("HttpPushPlugin.DestroyInstance", "传递的生产者为nil")
		return errors.New("传递的生产者为nil")
	}

	pushProducer, ok := prod.(*producer.Producer)
	if !ok {
		utils.PrintErr("HttpPushPlugin.DestroyInstance", "类型转换失败")
		return errors.New("类型转换失败")
	}

	producer.DestroyProducer(pushProducer)

	return nil
}
