package service

import (
	"com.fs/event-service/db"
)

var loader *pluginLoader

// Init 初始化service
func Init() {
	db.Open()

	loader = initPluginLoader()
	err := loader.load()
	if err != nil {
		panic(err)
	}
}

// Destroy 销毁service
func Destroy() {
	_ = loader.unload()
	destroyPluginLoader(loader)

	loader = nil

	db.Close()
}
