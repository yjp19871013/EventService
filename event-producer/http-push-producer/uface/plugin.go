package main

import (
	"com.fs/event-service/event-producer/http-push-producer/uface/producer"
)

var Factory producer.HttpPushFactory

func init() {
	Factory = producer.HttpPushFactory{}
	Factory.InitProducer = producer.InitProducer
	Factory.DestroyProducer = producer.DestroyProducer
}
