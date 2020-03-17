package main

import "com.fs/event-service/event-producer/http-push-producer/producer"

func main() {
	factory := &producer.HttpPushFactory{}
	prod, err := factory.NewInstance("http_push")
	if err != nil {
		panic(err)
	}

	err = factory.DestroyInstance(prod)
	if err != nil {
		panic(err)
	}
}
