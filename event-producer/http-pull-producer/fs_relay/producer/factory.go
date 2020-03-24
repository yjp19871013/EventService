package producer

import (
	"com.fs/event-service/event-producer/http-pull-producer/base"
	"com.fs/event-service/plugins"
	"com.fs/event-service/utils"
	"errors"
)

const (
	fsRelayInstancesDir = "fs_relay"
)

type HttpPullFactory struct {
	base.HttpPullFactory
}

func InitProducer(instanceName string, conf *base.Config) (plugins.Instance, error) {
	if conf == nil {
		utils.PrintErr("HttpPullFactory.InitProducer", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	prod := &HttpPullProducer{}
	prod.ProducerName = instanceName
	prod.Config = conf
	prod.Pull = prod.onPull

	return prod, nil
}

func DestroyProducer(instance plugins.Instance) error {
	if instance == nil {
		utils.PrintErr("HttpPullFactory.DestroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	pullProducer, ok := instance.(*HttpPullProducer)
	if !ok {
		utils.PrintErr("HttpPullFactory.DestroyProducer", "类型转换失败")
		return errors.New("类型转换失败")
	}

	pullProducer.Pull = nil
	pullProducer.Config = nil

	pullProducer = nil

	return nil
}

func OfferInstancesSubDir() string {
	return fsRelayInstancesDir
}
