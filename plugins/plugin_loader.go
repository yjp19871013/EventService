package plugins

import (
	"com.fs/event-service/config"
	"com.fs/event-service/utils"
	"errors"
	"fmt"
	"github.com/rjeczalik/notify"
	"plugin"
	"sync"
	"time"
)

type PluginLoader struct {
	pluginsDir string

	pluginMapLock sync.Mutex
	pluginMap     map[string]Plugin

	instanceLoaders []*instanceLoader

	stopChan chan bool
}

func InitPluginLoader(pluginsDir string) (*PluginLoader, error) {
	if utils.IsStringEmpty(pluginsDir) {
		utils.PrintErr("initPluginLoader", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	loader := &PluginLoader{}

	loader.pluginMapLock.Lock()
	loader.pluginMap = make(map[string]Plugin)
	loader.pluginMapLock.Unlock()

	loader.instanceLoaders = make([]*instanceLoader, 0)

	loader.stopChan = make(chan bool)

	loader.pluginsDir = pluginsDir

	return loader, nil
}

func DestroyPluginLoader(loader *PluginLoader) {
	if loader == nil {
		return
	}

	loader.pluginsDir = ""

	loader.stopChan = nil

	loader.pluginMapLock.Lock()
	loader.pluginMap = nil
	loader.pluginMapLock.Unlock()

	loader.instanceLoaders = nil

	loader = nil
}

func (loader *PluginLoader) Load() error {
	pluginFiles, err := loader.getAllPluginFiles()
	if err != nil {
		utils.PrintCallErr("pluginLoader.load", "loader.getAllPluginsNoLock", err)
		return err
	}

	for _, pluginFile := range pluginFiles {
		err := loader.loadPlugin(pluginFile)
		if err != nil {
			utils.PrintCallErr("pluginLoader.load", "loader.loadProducerPlugin", err)
			return err
		}
	}

	return nil
}

func (loader *PluginLoader) Unload() error {
	pluginFiles, err := loader.getAllPluginFiles()
	if err != nil {
		utils.PrintCallErr("pluginLoader.unload", "loader.getAllPluginFilesNoLock", err)
		return err
	}

	for _, pluginFile := range pluginFiles {
		err := loader.unloadPlugin(pluginFile)
		if err != nil {
			utils.PrintCallErr("pluginLoader.unload", "loader.unloadPlugin", err)
		}
	}

	return nil
}

func (loader *PluginLoader) Start() error {
	c := make(chan notify.EventInfo)

	err := notify.Watch(loader.pluginsDir, c, notify.Create, notify.Remove, notify.Rename)
	if err != nil {
		utils.PrintCallErr("pluginLoader.start", "notify.Watch", err)
		return err
	}

	go func(loader *PluginLoader, c chan notify.EventInfo) {
		for true {
			select {
			case stop := <-loader.stopChan:
				if stop {
					fmt.Println("PluginLoader Stopped")
					notify.Stop(c)
					loader.stopChan <- true
				}
			case ei := <-c:
				switch ei.Event() {
				case notify.Create:
					go func() {
						var i uint
						conf := config.GetEventServiceConfig().PluginConfig
						for i = 0; i < conf.LoadPluginTryTimes; i++ {
							err = loader.loadPlugin(ei.Path())
							if err != nil {
								time.Sleep(time.Duration(conf.LoadPluginTimeoutSec) * time.Second)
								continue
							}

							break
						}

						if i == conf.LoadPluginTryTimes {
							fmt.Println(ei.Path(), "loaded failed:", err.Error())
						} else {
							fmt.Println(ei.Path(), "loaded")
						}
					}()
				case notify.Remove:
					fallthrough
				case notify.Rename:
					err := loader.unloadPlugin(ei.Path())
					if err != nil {
						utils.PrintCallErr("PluginLoader.Start", "loader.unloadPlugin", err)
						continue
					}

					fmt.Println(ei.Path() + " unloaded")
				}
			}
		}
	}(loader, c)

	return nil
}

func (loader *PluginLoader) Stop() {
	loader.stopChan <- true
	<-loader.stopChan
}

func (loader *PluginLoader) loadPlugin(pluginFilePath string) error {
	if utils.IsStringEmpty(pluginFilePath) {
		utils.PrintErr("loadPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	exist := utils.PathExists(pluginFilePath)
	if !exist {
		utils.PrintErr("loadPlugin", pluginFilePath+"插件不存在")
		return errors.New(pluginFilePath + "插件不存在")
	}

	p, err := plugin.Open(pluginFilePath)
	if err != nil {
		utils.PrintCallErr("loadPlugin", "plugin.Open", err)
		return err
	}

	s, err := p.Lookup("Plugin")
	if err != nil {
		utils.PrintCallErr("loadPlugin", "p.Lookup", err)
		return err
	}

	instancePlugin, ok := s.(Plugin)
	if !ok {
		utils.PrintErr("loadPlugin", "类型转换失败")
		return errors.New("类型转换失败")
	}

	loader.addPlugin(pluginFilePath, instancePlugin)

	instanceLoader, err := initInstanceLoader(instancePlugin)
	if err != nil {
		utils.PrintCallErr("loadPlugin", "initInstanceLoader", err)
		return err
	}

	err = instanceLoader.start()
	if err != nil {
		utils.PrintCallErr("loadPlugin", "instanceLoader.start", err)
		return err
	}

	loader.instanceLoaders = append(loader.instanceLoaders, instanceLoader)

	return nil
}

func (loader *PluginLoader) unloadPlugin(pluginFilePath string) error {
	if utils.IsStringEmpty(pluginFilePath) {
		utils.PrintErr("pluginLoader.unloadPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	for _, instanceLoader := range loader.instanceLoaders {
		instanceLoader.stop()
		destroyInstanceLoader(instanceLoader)
	}

	loader.deletePlugin(pluginFilePath)

	return nil
}

func (loader *PluginLoader) addPlugin(pluginFilePath string, p Plugin) {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	loader.pluginMap[pluginFilePath] = p
}

func (loader *PluginLoader) deletePlugin(pluginFilePath string) {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	loader.pluginMap[pluginFilePath] = nil
	delete(loader.pluginMap, pluginFilePath)
}

func (loader *PluginLoader) getPlugin(pluginFilePath string) Plugin {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	return loader.pluginMap[pluginFilePath]
}

func (loader *PluginLoader) getAllPluginFiles() ([]string, error) {
	return utils.GetDirFiles(loader.pluginsDir)
}

func (loader *PluginLoader) getAllLoadedPlugins() []string {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	plugins := make([]string, 0)
	for pluginFilePath, _ := range loader.pluginMap {
		plugins = append(plugins, pluginFilePath)
	}

	return plugins
}
