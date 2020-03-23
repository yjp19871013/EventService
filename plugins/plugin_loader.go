package plugins

import (
	"com.fs/event-service/config"
	"com.fs/event-service/utils"
	"errors"
	"fmt"
	"github.com/rjeczalik/notify"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"time"
)

type pluginLoader struct {
	pluginsDir   string
	instancesDir string

	pluginMapLock sync.Mutex
	pluginMap     map[string]Plugin

	instanceLoaderMapLock sync.Mutex
	instanceLoaderMap     map[string]instanceLoader

	stopChan chan bool
}

func InitPluginLoader(pluginsDir string, instancesDir string) (*pluginLoader, error) {
	if utils.IsStringEmpty(pluginsDir) || utils.IsStringEmpty(instancesDir) {
		utils.PrintErr("initPluginLoader", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	loader := &pluginLoader{}

	loader.pluginMapLock.Lock()
	loader.pluginMap = make(map[string]Plugin)
	loader.pluginMapLock.Unlock()

	loader.instanceLoaderMapLock.Lock()
	loader.instanceLoaderMap = make(map[string]instanceLoader)
	loader.instanceLoaderMapLock.Unlock()

	loader.stopChan = make(chan bool)

	loader.pluginsDir = pluginsDir
	loader.instancesDir = instancesDir

	return loader, nil
}

func DestroyPluginLoader(loader *pluginLoader) {
	if loader == nil {
		return
	}

	loader.instancesDir = ""
	loader.pluginsDir = ""

	loader.stopChan = nil

	loader.pluginMapLock.Lock()
	loader.pluginMap = nil
	loader.pluginMapLock.Unlock()

	loader.instanceLoaderMapLock.Lock()
	loader.instanceLoaderMap = nil
	loader.instanceLoaderMapLock.Unlock()

	loader = nil
}

func (loader *pluginLoader) Load() error {
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

func (loader *pluginLoader) Unload() error {
	pluginFiles, err := loader.getAllPluginFiles()
	if err != nil {
		utils.PrintCallErr("pluginLoader.unload", "loader.getAllPluginFilesNoLock", err)
		return err
	}

	for _, pluginFile := range pluginFiles {
		loader.unloadPlugin(pluginFile)
	}

	return nil
}

func (loader *pluginLoader) Start() error {
	c := make(chan notify.EventInfo)

	err := notify.Watch(loader.pluginsDir, c, notify.Create, notify.Remove, notify.Rename)
	if err != nil {
		utils.PrintCallErr("pluginLoader.start", "notify.Watch", err)
		return err
	}

	go func(loader *pluginLoader, c chan notify.EventInfo) {
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
					loader.unloadPlugin(ei.Path())
					fmt.Println(ei.Path() + " unloaded")
				}
			}
		}
	}(loader, c)

	return nil
}

func (loader *pluginLoader) Stop() {
	loader.stopChan <- true
	<-loader.stopChan
}

func (loader *pluginLoader) loadPlugin(pluginFilePath string) error {
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

	pluginFileName := filepath.Base(pluginFilePath)
	instancesSubDirName := strings.ReplaceAll(pluginFileName, filepath.Ext(pluginFileName), "")
	instancesSubDir := filepath.Join(loader.instancesDir, instancesSubDirName)
	if !utils.PathExists(instancesSubDir) {
		err := os.MkdirAll(instancesSubDir, os.ModeDir|os.ModePerm)
		if err != nil {
			utils.PrintCallErr("loadPlugin", "os.MkdirAll", err)
			return err
		}
	}

	instanceLoader, err := initInstanceLoader(instancePlugin, instancesSubDir)
	if err != nil {
		utils.PrintCallErr("loadPlugin", "initInstanceLoader", err)
		return err
	}

	err = instanceLoader.start()
	if err != nil {
		utils.PrintCallErr("loadPlugin", "instanceLoader.start", err)
		return err
	}

	loader.addInstanceLoader(pluginFilePath, instanceLoader)

	return nil
}

func (loader *pluginLoader) unloadPlugin(pluginFilePath string) error {
	if utils.IsStringEmpty(pluginFilePath) {
		utils.PrintErr("pluginLoader.unloadPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	for _, instanceLoader := range loader.instanceLoaderMap {
		instanceLoader.stop()
		destroyInstanceLoader(&instanceLoader)
	}

	loader.deletePlugin(pluginFilePath)

	return nil
}

func (loader *pluginLoader) addPlugin(pluginFilePath string, p Plugin) {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	loader.pluginMap[pluginFilePath] = p
}

func (loader *pluginLoader) deletePlugin(pluginFilePath string) {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	loader.pluginMap[pluginFilePath] = nil
	delete(loader.pluginMap, pluginFilePath)
}

func (loader *pluginLoader) getPlugin(pluginFilePath string) Plugin {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	return loader.pluginMap[pluginFilePath]
}

func (loader *pluginLoader) addInstanceLoader(pluginFilePath string, instanceLoader *instanceLoader) {
	loader.instanceLoaderMapLock.Lock()
	defer loader.instanceLoaderMapLock.Unlock()

	loader.instanceLoaderMap[pluginFilePath] = *instanceLoader
}

func (loader *pluginLoader) getAllPluginFiles() ([]string, error) {
	return utils.GetDirFiles(loader.pluginsDir)
}

func (loader *pluginLoader) getAllLoadedPlugins() []string {
	loader.pluginMapLock.Lock()
	defer loader.pluginMapLock.Unlock()

	plugins := make([]string, 0)
	for pluginFilePath, _ := range loader.pluginMap {
		plugins = append(plugins, pluginFilePath)
	}

	return plugins
}
