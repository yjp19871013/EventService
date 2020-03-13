package main

import (
	"com.fs/event-service/event-producer"
	"fmt"
)

func init() {
	fmt.Println("http push producer load")
}

type HttpPushPlugin struct {
}

var Plugin HttpPushPlugin

func (p *HttpPushPlugin) NewInstance(conf string) event_producer.EventProducer {
	fmt.Println("httpPushPlugin NewInstance")

	return nil
}

func (p *HttpPushPlugin) DestroyInstance(producer event_producer.EventProducer) {
	if producer == nil {
		return
	}

	fmt.Println("httpPushPlugin DestroyInstance")
}
