package event_producer

type EventProducerPlugin interface {
	NewInstance(conf string) (EventProducer, error)
	DestroyInstance(producer EventProducer) error
}

type EventProducer interface {
}
