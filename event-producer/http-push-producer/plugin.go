package main

import (
	"com.fs/event-service/event-producer"
	"com.fs/event-service/event-producer/http-push-producer/producer"
	"fmt"
)

func init() {
	fmt.Println("http push producer load")
}

type HttpPushPlugin struct {
}

var Plugin = HttpPushPlugin{}

func (p *HttpPushPlugin) NewInstance(conf string) event_producer.EventProducer {
	fmt.Println("httpPushPlugin NewInstance")

	return &producer.Producer{}
}

func (p *HttpPushPlugin) DestroyInstance(producer event_producer.EventProducer) {
	if producer == nil {
		return
	}

	fmt.Println("httpPushPlugin DestroyInstance")
}
