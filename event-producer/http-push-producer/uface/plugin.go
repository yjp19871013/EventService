package main

import (
	"com.fs/event-service/event-producer/http-push-producer/uface/producer"
)

var Plugin producer.HttpPushFactory

func init() {
	Plugin = producer.HttpPushFactory{}
	Plugin.InitProducer = producer.InitProducer
	Plugin.DestroyProducer = producer.DestroyProducer
}
