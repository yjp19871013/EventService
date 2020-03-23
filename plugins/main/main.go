package main

import "com.fs/event-service/plugins"

func main() {
	pluginLoader, err := plugins.InitPluginLoader("/home/yjp/go-projects/EventService/src/com.fs/event-service/deployment/plugins")
	if err != nil {
		panic(err)
	}
	defer plugins.DestroyPluginLoader(pluginLoader)

	pluginLoader.Load()
	defer pluginLoader.Unload()

	pluginLoader.Start()
	defer pluginLoader.Stop()

	for true {
		continue
	}
}
