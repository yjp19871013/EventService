package service

import (
	"com.fs/event-service/config"
	"com.fs/event-service/db"
	"com.fs/event-service/event-producer"
	"com.fs/event-service/utils"
	"errors"
	"path/filepath"
	"plugin"
	"sync"
)

type pluginLoader struct {
	pluginMapLock sync.Mutex
	pluginMap     map[string]event_producer.EventProducerFactory

	producerMapLock sync.Mutex
	producerMap     map[string]event_producer.EventProducer
}

func initPluginLoader() *pluginLoader {
	loader := &pluginLoader{}

	loader.pluginMap = make(map[string]event_producer.EventProducerFactory)
	loader.producerMap = make(map[string]event_producer.EventProducer)

	return loader
}

func destroyPluginLoader(loader *pluginLoader) {
	if loader == nil {
		return
	}

	loader.producerMap = nil
	loader.pluginMap = nil

	loader = nil
}

func (loader *pluginLoader) load() error {
	ps, err := db.GetAllProducerPlugins()
	if err != nil {
		utils.PrintCallErr("pluginLoader.startLoader", "db.GetAllProducerPlugins", err)
		return err
	}

	for _, p := range ps {
		err := loader.loadProducerPlugin(p.PluginFileName)
		if err != nil {
			utils.PrintCallErr("pluginLoader.startLoader", "loader.loadProducerPlugin", err)
			return err
		}

		producers, err := p.GetAllPluginProducers()
		if err != nil {
			utils.PrintCallErr("pluginLoader.startLoader", "p.GetAllPluginProducers", err)
			return err
		}

		for _, producer := range producers {
			err := loader.newProducer(p.PluginFileName, producer.Name)
			if err != nil {
				utils.PrintCallErr("pluginLoader.startLoader", "loader.newProducer", err)
				return err
			}
		}
	}

	return nil
}

func (loader *pluginLoader) unload() error {
	ps, err := db.GetAllProducerPlugins()
	if err != nil {
		utils.PrintCallErr("pluginLoader.stopLoader", "db.GetAllProducerPlugins", err)
		return err
	}

	for _, p := range ps {
		producers, err := p.GetAllPluginProducers()
		if err != nil {
			utils.PrintCallErr("pluginLoader.stopLoader", "p.GetAllPluginProducers", err)
			return err
		}

		for _, producer := range producers {
			err := loader.destroyProducer(p.PluginFileName, producer.Name)
			if err != nil {
				utils.PrintCallErr("pluginLoader.stopLoader", "loader.destroyProducer", err)
				return err
			}
		}

		loader.unloadProducerPlugin(p.PluginFileName)
	}

	return nil
}

func (loader *pluginLoader) loadProducerPlugin(pluginFileName string) error {
	if utils.IsStringEmpty(pluginFileName) {
		utils.PrintErr("loadProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	conf := config.GetEventServiceConfig().PluginConfig
	pluginPath := filepath.Join(conf.Dir, pluginFileName)
	exist := utils.PathExists(pluginPath)
	if !exist {
		utils.PrintErr("loadProducerPlugin", pluginFileName+"插件不存在")
		return errors.New(pluginFileName + "插件不存在")
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		utils.PrintCallErr("loadProducerPlugin", "plugin.Open", err)
		return err
	}

	s, err := p.Lookup("Factory")
	if err != nil {
		utils.PrintCallErr("newProducer", "p.Lookup", err)
		return err
	}

	producerPlugin, ok := s.(event_producer.EventProducerFactory)
	if !ok {
		utils.PrintErr("newProducer", "类型转换失败")
		return errors.New("类型转换失败")
	}

	loader.addProducerPlugin(pluginFileName, producerPlugin)

	return nil
}

func (loader *pluginLoader) unloadProducerPlugin(pluginFileName string) {
	loader.deleteProducerPlugin(pluginFileName)
}

func (loader *pluginLoader) newProducer(pluginFileName string, producerName string) error {
	if utils.IsStringEmpty(pluginFileName) || utils.IsStringEmpty(producerName) {
		utils.PrintErr("newProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := loader.getProducer(producerName)
	if producer != nil {
		utils.PrintErr("newProducer", "该名称的生产者已存在")
		return errors.New("该名称的生产者已存在")
	}

	p := loader.getProducerPlugin(pluginFileName)
	if p == nil {
		utils.PrintErr("newProducer", "没有找到对应的插件")
		return errors.New("没有找到对应的插件")
	}

	producer, err := p.NewInstance(producerName)
	if err != nil {
		utils.PrintCallErr("newProducer", "p.NewInstance", err)
		return err
	}

	loader.addProducer(producerName, producer)

	return nil
}

func (loader *pluginLoader) destroyProducer(pluginFileName string, producerName string) error {
	if utils.IsStringEmpty(pluginFileName) || utils.IsStringEmpty(producerName) {
		utils.PrintErr("destroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := loader.getProducerPlugin(pluginFileName)
	if p == nil {
		utils.PrintErr("destroyProducer", "没有找到对应的插件")
		return errors.New("没有找到对应的插件")
	}

	producer := loader.getProducer(producerName)
	if producer == nil {
		utils.PrintErr("destroyProducer", "没有找到对应的生产者")
		return errors.New("没有找到对应的生产者")
	}

	loader.deleteProducer(producerName)

	err := p.DestroyInstance(producer)
	if err != nil {
		utils.PrintCallErr("destroyProducer", "p.DestroyInstance", err)
		return err
	}

	return nil
}

func (loader *pluginLoader) addProducerPlugin(pluginFileName string, p event_producer.EventProducerFactory) {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	loader.pluginMap[pluginFileName] = p
}

func (loader *pluginLoader) deleteProducerPlugin(pluginFileName string) {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	loader.pluginMap[pluginFileName] = nil
	delete(loader.pluginMap, pluginFileName)
}

func (loader *pluginLoader) getProducerPlugin(pluginFileName string) event_producer.EventProducerFactory {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	return loader.pluginMap[pluginFileName]
}

func (loader *pluginLoader) getAllProducerPlugins() []string {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	plugins := make([]string, 0)
	for pluginFileName, _ := range loader.pluginMap {
		plugins = append(plugins, pluginFileName)
	}

	return plugins
}

func (loader *pluginLoader) addProducer(producerName string, producer event_producer.EventProducer) {
	loader.producerMapLock.Lock()
	defer loader.producerMapLock.Unlock()

	loader.producerMap[producerName] = producer
}

func (loader *pluginLoader) getProducer(producerName string) event_producer.EventProducer {
	loader.producerMapLock.Lock()
	defer loader.producerMapLock.Unlock()

	return loader.producerMap[producerName]
}

func (loader *pluginLoader) deleteProducer(producerName string) {
	loader.producerMapLock.Lock()
	defer loader.producerMapLock.Unlock()

	loader.producerMap[producerName] = nil
	delete(loader.producerMap, producerName)
}

func (loader *pluginLoader) getAllProducers() []string {
	loader.producerMapLock.Lock()
	defer loader.producerMapLock.Unlock()

	producerNames := make([]string, 0)
	for producerName, _ := range loader.producerMap {
		producerNames = append(producerNames, producerName)
	}

	return producerNames
}
