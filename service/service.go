package service

import (
	"com.fs/event-service/db"
)

// Init 初始化service
func Init() {
	db.Open()
	initPluginsAndProducers()
}

// Destroy 销毁service
func Destroy() {
	destroyPluginsAndProducers()
	db.Close()
}

func initPluginsAndProducers() {
	ps, err := db.GetAllProducerPlugins()
	if err != nil {
		panic(err)
	}

	for _, p := range ps {
		err := loadProducerPlugin(p.Name)
		if err != nil {
			panic("Load: " + p.Name + err.Error())
		}

		producers, err := p.GetAllPluginProducers()
		if err != nil {
			panic(p.Name + " GetAllPluginProducers: " + err.Error())
		}

		for _, producer := range producers {
			err := newProducer(p.Name, producer.Name, producer.Config)
			if err != nil {
				panic(p.Name + "=" + producer.Name + " newProducer: " + err.Error())
			}
		}
	}
}

func destroyPluginsAndProducers() {
	ps, err := db.GetAllProducerPlugins()
	if err != nil {
		panic(err)
	}

	for _, p := range ps {
		producers, err := p.GetAllPluginProducers()
		if err != nil {
			panic(p.Name + " GetAllPluginProducers: " + err.Error())
		}

		for _, producer := range producers {
			err := destroyProducer(p.Name, producer.Name)
			if err != nil {
				panic(p.Name + "=" + producer.Name + " destroyProducer: " + err.Error())
			}
		}

		unloadProducerPlugin(p.Name)
	}
}
