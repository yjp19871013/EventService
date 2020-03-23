package plugins

//import (
//	"com.fs/event-service/config"
//	"com.fs/event-service/db"
//	"com.fs/event-service/event-producer"
//	"com.fs/event-service/utils"
//	"errors"
//	"path/filepath"
//	"plugin"
//)
//
//func (loader *pluginLoader) newProducer(pluginFileName string, producerName string) error {
//	if utils.IsStringEmpty(pluginFileName) || utils.IsStringEmpty(producerName) {
//		utils.PrintErr("newProducer", "没有传递必要的参数")
//		return errors.New("没有传递必要的参数")
//	}
//
//	producer := loader.getProducer(producerName)
//	if producer != nil {
//		utils.PrintErr("newProducer", "该名称的生产者已存在")
//		return errors.New("该名称的生产者已存在")
//	}
//
//	p := loader.getProducerPlugin(pluginFileName)
//	if p == nil {
//		utils.PrintErr("newProducer", "没有找到对应的插件")
//		return errors.New("没有找到对应的插件")
//	}
//
//	producer, err := p.NewInstance(producerName)
//	if err != nil {
//		utils.PrintCallErr("newProducer", "p.NewInstance", err)
//		return err
//	}
//
//	err = producer.Start()
//	if err != nil {
//		utils.PrintCallErr("newProducer", "producer.Start", err)
//		return err
//	}
//
//	loader.addProducer(producerName, producer)
//
//	return nil
//}
//
//func (loader *pluginLoader) destroyProducer(pluginFileName string, producerName string) error {
//	if utils.IsStringEmpty(pluginFileName) || utils.IsStringEmpty(producerName) {
//		utils.PrintErr("destroyProducer", "没有传递必要的参数")
//		return errors.New("没有传递必要的参数")
//	}
//
//	p := loader.getProducerPlugin(pluginFileName)
//	if p == nil {
//		utils.PrintErr("destroyProducer", "没有找到对应的插件")
//		return errors.New("没有找到对应的插件")
//	}
//
//	producer := loader.getProducer(producerName)
//	if producer == nil {
//		utils.PrintErr("destroyProducer", "没有找到对应的生产者")
//		return errors.New("没有找到对应的生产者")
//	}
//
//	loader.deleteProducer(producerName)
//
//	err := producer.Stop()
//	if err != nil {
//		utils.PrintCallErr("destroyProducer", "producer.Stop", err)
//		return err
//	}
//
//	err = p.DestroyInstance(producer)
//	if err != nil {
//		utils.PrintCallErr("destroyProducer", "p.DestroyInstance", err)
//		return err
//	}
//
//	return nil
//}
//
//func (loader *pluginLoader) addProducer(producerName string, producer event_producer.EventProducer) {
//	loader.producerMapLock.Lock()
//	defer loader.producerMapLock.Unlock()
//
//	loader.producerMap[producerName] = producer
//}
//
//func (loader *pluginLoader) getProducer(producerName string) event_producer.EventProducer {
//	loader.producerMapLock.Lock()
//	defer loader.producerMapLock.Unlock()
//
//	return loader.producerMap[producerName]
//}
//
//func (loader *pluginLoader) deleteProducer(producerName string) {
//	loader.producerMapLock.Lock()
//	defer loader.producerMapLock.Unlock()
//
//	loader.producerMap[producerName] = nil
//	delete(loader.producerMap, producerName)
//}
//
//func (loader *pluginLoader) getAllProducers() []string {
//	loader.producerMapLock.Lock()
//	defer loader.producerMapLock.Unlock()
//
//	producerNames := make([]string, 0)
//	for producerName, _ := range loader.producerMap {
//		producerNames = append(producerNames, producerName)
//	}
//
//	return producerNames
//}
