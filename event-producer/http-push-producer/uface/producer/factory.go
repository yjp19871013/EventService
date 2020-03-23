package producer

import (
	"com.fs/event-service/event-producer/http-push-producer/base"
	"com.fs/event-service/plugins"
	"com.fs/event-service/utils"
	"errors"
)

type HttpPushFactory struct {
	base.HttpPushFactory
}

func InitProducer(conf *base.Config) (plugins.Instance, error) {
	if conf == nil {
		utils.PrintErr("HttpPushFactory.InitProducer", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	prod := &HttpPushProducer{}
	prod.Config = conf
	prod.OnHandle = prod.handlerFun

	return prod, nil
}

func DestroyProducer(instance plugins.Instance) error {
	if instance == nil {
		utils.PrintErr("HttpPushFactory.DestroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	pushProducer, ok := instance.(*HttpPushProducer)
	if !ok {
		utils.PrintErr("HttpPushFactory.DestroyProducer", "类型转换失败")
		return errors.New("类型转换失败")
	}

	pushProducer.OnHandle = nil
	pushProducer.Config = nil

	pushProducer = nil

	return nil
}
