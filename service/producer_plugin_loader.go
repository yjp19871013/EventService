package service

import (
	"com.fs/event-service/config"
	"com.fs/event-service/event-producer"
	"com.fs/event-service/utils"
	"errors"
	"path/filepath"
	"plugin"
	"sync"
)

const (
	pluginExt = ".so"
)

var pluginMapLock sync.Mutex
var pluginMap = map[string]event_producer.EventProducerPlugin{}

func loadProducerPlugin(pluginName string) error {
	if utils.IsStringEmpty(pluginName) {
		utils.PrintErr("loadProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	conf := config.GetEventServiceConfig().PluginConfig
	pluginPath := filepath.Join(conf.Dir, pluginName+pluginExt)
	exist := utils.PathExists(pluginPath)
	if !exist {
		utils.PrintErr("loadProducerPlugin", pluginName+"插件不存在")
		return errors.New(pluginName + "插件不存在")
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		utils.PrintCallErr("loadProducerPlugin", "plugin.Open", err)
		return err
	}

	s, err := p.Lookup("Plugin")
	if err != nil {
		utils.PrintCallErr("newProducer", "p.Lookup", err)
		return err
	}

	producerPlugin, ok := s.(event_producer.EventProducerPlugin)
	if !ok {
		utils.PrintErr("newProducer", "类型转换失败")
		return errors.New("类型转换失败")
	}

	addProducerPlugin(pluginName, producerPlugin)

	return nil
}

func unloadProducerPlugin(pluginName string) {
	deleteProducerPlugin(pluginName)
}

func addProducerPlugin(pluginName string, p event_producer.EventProducerPlugin) {
	pluginMapLock.Lock()
	defer pluginMapLock.Unlock()

	pluginMap[pluginName] = p
}

func deleteProducerPlugin(pluginName string) {
	pluginMapLock.Lock()
	defer pluginMapLock.Unlock()

	pluginMap[pluginName] = nil
	delete(pluginMap, pluginName)
}

func getProducerPlugin(pluginName string) event_producer.EventProducerPlugin {
	pluginMapLock.Lock()
	defer pluginMapLock.Unlock()

	return pluginMap[pluginName]
}
