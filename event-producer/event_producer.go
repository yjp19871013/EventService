package event_producer

type EventProducerPlugin interface {
	NewInstance(conf string) EventProducer
	DestroyInstance(producer EventProducer)
}

type EventProducer interface {
}
