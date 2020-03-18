package main

import "com.fs/event-service/event-producer/http-pull-producer/fs_relay/producer"

var Factory producer.HttpPullFactory

func init() {
	Factory = producer.HttpPullFactory{}
	Factory.InitProducer = producer.InitProducer
	Factory.DestroyProducer = producer.DestroyProducer
}
