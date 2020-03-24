package service

import (
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"errors"
)

func GetAllProducers() []model.ProducerInfo {
	return model.TransferInstancesToProducerInfoBatch(loader.GetAllInstances())
}

func GetPluginProducers(pluginName string) (*model.ProducerInfo, error) {
	if utils.IsStringEmpty(pluginName) {
		utils.PrintErr("GetPluginProducers", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	instances, err := loader.GetPluginInstances(pluginName)
	if err != nil {
		utils.PrintCallErr("GetPluginProducers", "loader.GetPluginInstances", err)
		return nil, err
	}

	return model.TransferInstancesToProducerInfo(pluginName, instances), nil
}
