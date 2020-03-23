package event_producer

import (
	"com.fs/event-service/config"
	"com.fs/event-service/plugins"
	"com.fs/event-service/utils"
	"os"
	"path/filepath"
)

func init() {
	conf := config.GetEventServiceConfig().PluginConfig
	pluginsDir := filepath.Join(conf.Dir)
	if !utils.PathExists(pluginsDir) {
		err := os.MkdirAll(pluginsDir, os.ModePerm|os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
}

type EventProducerFactory interface {
	plugins.Plugin
}

type EventProducer interface {
	plugins.Instance
}
