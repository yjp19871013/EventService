package main

import "com.fs/event-service/event-producer/http-pull-producer/fs_relay/producer"

var Plugin producer.HttpPullFactory

func init() {
	Plugin = producer.HttpPullFactory{}
	Plugin.InitProducer = producer.InitProducer
	Plugin.DestroyProducer = producer.DestroyProducer
}
