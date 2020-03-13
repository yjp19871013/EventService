package service

import (
	"com.fs/event-service/db"
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"errors"
)

func AddProducerPlugin(pluginName string) error {
	if utils.IsStringEmpty(pluginName) {
		utils.PrintErr("AddProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{Name: pluginName}
	err := p.Create()
	if err != nil {
		utils.PrintCallErr("AddProducerPlugin", "p.Create", err)
		return err
	}

	return nil
}

func DeleteProducerPlugin(pluginName string) error {
	if utils.IsStringEmpty(pluginName) {
		utils.PrintErr("DeleteProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{Name: pluginName}
	err := p.GetByName()
	if err != nil {
		utils.PrintCallErr("DeleteProducerPlugin", "p.GetByName", err)
		return err
	}

	err = p.DeleteByIDAndName()
	if err != nil {
		utils.PrintCallErr("DeleteProducerPlugin", "p.DeleteByIDAndName", err)
		return err
	}

	return nil
}

func GetProducerPlugins() ([]model.ProducerPluginInfo, error) {
	plugins, err := db.GetAllProducerPlugins()
	if err != nil {
		return nil, err
	}

	return model.TransferProducerPluginToProducerPluginInfoBatch(plugins), nil
}
