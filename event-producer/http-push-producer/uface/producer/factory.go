package producer

import (
	"com.fs/event-service/event-producer"
	"com.fs/event-service/event-producer/http-push-producer/base"
	"com.fs/event-service/utils"
	"errors"
)

type HttpPushFactory struct {
	base.HttpPushFactory
}

func InitProducer(conf *base.Config) (event_producer.EventProducer, error) {
	if conf == nil {
		utils.PrintErr("HttpPushFactory.InitProducer", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	prod := &HttpPushProducer{}
	prod.Config = conf
	prod.OnHandle = HandlerFun

	return prod, nil
}

func DestroyProducer(prod event_producer.EventProducer) error {
	if prod == nil {
		utils.PrintErr("HttpPushFactory.DestroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	pushProducer, ok := prod.(*HttpPushProducer)
	if !ok {
		utils.PrintErr("HttpPushFactory.DestroyProducer", "类型转换失败")
		return errors.New("类型转换失败")
	}

	pushProducer.Config = nil
	pushProducer = nil

	return nil
}
