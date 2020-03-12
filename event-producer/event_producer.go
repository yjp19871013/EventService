package event_producer

type EventProducerFactory interface {
	NewInstance(conf string)
	DestroyInstance(producer EventProducer)
}

type EventProducer interface {
}
