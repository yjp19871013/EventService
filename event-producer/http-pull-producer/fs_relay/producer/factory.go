package producer

import (
	"com.fs/event-service/event-producer"
	"com.fs/event-service/event-producer/http-pull-producer/base"
	"com.fs/event-service/utils"
	"errors"
)

type HttpPullFactory struct {
	base.HttpPullFactory
}

func InitProducer(producerName string, conf *base.Config) (event_producer.EventProducer, error) {
	if utils.IsStringEmpty(producerName) || conf == nil {
		utils.PrintErr("HttpPullFactory.InitProducer", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	prod := &HttpPullProducer{}
	prod.ProducerName = producerName
	prod.Config = conf
	prod.Pull = prod.onPull

	return prod, nil
}

func DestroyProducer(prod event_producer.EventProducer) error {
	if prod == nil {
		utils.PrintErr("HttpPullFactory.DestroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	pullProducer, ok := prod.(*HttpPullProducer)
	if !ok {
		utils.PrintErr("HttpPullFactory.DestroyProducer", "类型转换失败")
		return errors.New("类型转换失败")
	}

	pullProducer.Pull = nil
	pullProducer.Config = nil
	pullProducer.ProducerName = ""

	pullProducer = nil

	return nil
}
