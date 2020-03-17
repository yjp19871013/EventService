package producer

import (
	"com.fs/event-service/event-producer"
	"com.fs/event-service/utils"
	"errors"
)

type HttpPushFactory struct {
}

func (p *HttpPushFactory) NewInstance(producerName string) (event_producer.EventProducer, error) {
	pushProducer, err := InitProducer(producerName)
	if err != nil {
		utils.PrintCallErr("HttpPushFactory.NewInstance", "producer.InitProducer", err)
		return nil, err
	}

	return pushProducer, nil
}

func (p *HttpPushFactory) DestroyInstance(prod event_producer.EventProducer) error {
	if prod == nil {
		utils.PrintErr("HttpPushFactory.DestroyInstance", "传递的生产者为nil")
		return errors.New("传递的生产者为nil")
	}

	pushProducer, ok := prod.(*Producer)
	if !ok {
		utils.PrintErr("HttpPushFactory.DestroyInstance", "类型转换失败")
		return errors.New("类型转换失败")
	}

	DestroyProducer(pushProducer)

	return nil
}
