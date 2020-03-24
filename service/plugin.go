package service

func GetLoadedPlugins() []string {
	return loader.GetAllLoadedPlugins()
}
