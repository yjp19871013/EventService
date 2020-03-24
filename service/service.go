package service

import (
	"com.fs/event-service/config"
	"com.fs/event-service/db"
	"com.fs/event-service/plugins"
	"com.fs/event-service/utils"
)

var loader *plugins.PluginLoader

// Init 初始化service
func Init() {
	db.Open()

	err := initPluginLoader()
	if err != nil {
		panic(err)
	}
}

// Destroy 销毁service
func Destroy() {
	destroyPluginLoader()

	db.Close()
}

func initPluginLoader() error {
	pluginConfig := config.GetEventServiceConfig().PluginConfig
	pluginLoader, err := plugins.InitPluginLoader(pluginConfig.Dir)
	if err != nil {
		utils.PrintCallErr("initPluginLoader", "plugins.InitPluginLoader", err)
		return err
	}

	loader = pluginLoader

	err = loader.Load()
	if err != nil {
		utils.PrintCallErr("initPluginLoader", "loader.Load", err)
		return err
	}

	err = loader.Start()
	if err != nil {
		utils.PrintCallErr("initPluginLoader", "loader.Start", err)
		return err
	}

	return nil
}

func destroyPluginLoader() {
	loader.Stop()
	_ = loader.Unload()

	plugins.DestroyPluginLoader(loader)

	loader = nil
}
