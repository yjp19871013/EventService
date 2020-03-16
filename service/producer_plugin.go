package service

import (
	"com.fs/event-service/db"
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"errors"
)

func AddProducerPlugin(pluginName string, pluginFileName string) error {
	if utils.IsStringEmpty(pluginName) || utils.IsStringEmpty(pluginFileName) {
		utils.PrintErr("AddProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{Name: pluginName, PluginFileName: pluginFileName}
	err := p.Create()
	if err != nil {
		utils.PrintCallErr("AddProducerPlugin", "p.Create", err)
		return err
	}

	return nil
}

func DeleteProducerPlugin(pluginID uint64) error {
	if pluginID == 0 {
		utils.PrintErr("DeleteProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{ID: pluginID}
	err := p.GetByID()
	if err != nil {
		utils.PrintCallErr("DeleteProducerPlugin", "p.GetByID", err)
		return err
	}

	err = p.DeleteByID()
	if err != nil {
		utils.PrintCallErr("DeleteProducerPlugin", "p.DeleteByID", err)
		return err
	}

	return nil
}

func GetProducerPlugins() ([]model.ProducerPluginInfo, error) {
	plugins, err := db.GetAllProducerPlugins()
	if err != nil {
		utils.PrintCallErr("GetProducerPlugins", "db.GetAllProducerPlugins", err)
		return nil, err
	}

	return model.TransferProducerPluginToProducerPluginInfoBatch(plugins), nil
}
