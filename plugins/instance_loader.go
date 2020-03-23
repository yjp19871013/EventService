package plugins

import (
	"com.fs/event-service/config"
	"com.fs/event-service/utils"
	"errors"
	"fmt"
	"github.com/rjeczalik/notify"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type instanceLoader struct {
	plugin       Plugin
	instancesDir string

	instanceMapLock sync.Mutex
	instanceMap     map[string]Instance

	stopChan chan bool
}

func initInstanceLoader(plugin Plugin, instancesDir string) (*instanceLoader, error) {
	if plugin == nil || utils.IsStringEmpty(instancesDir) {
		utils.PrintErr("initInstanceLoader", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	loader := &instanceLoader{}

	loader.plugin = plugin
	loader.instancesDir = instancesDir

	loader.instanceMapLock.Lock()
	loader.instanceMap = make(map[string]Instance)
	loader.instanceMapLock.Unlock()

	loader.stopChan = make(chan bool)

	return loader, nil
}

func destroyInstanceLoader(loader *instanceLoader) {
	if loader == nil {
		return
	}

	loader.stopChan = nil

	loader.instanceMapLock.Lock()
	loader.instanceMap = nil
	loader.instanceMapLock.Unlock()

	loader.instancesDir = ""
	loader.plugin = nil

	loader = nil
}

func (loader *instanceLoader) start() error {
	c := make(chan notify.EventInfo)

	err := notify.Watch(loader.instancesDir, c, notify.Create, notify.Remove, notify.Rename)
	if err != nil {
		utils.PrintCallErr("instanceLoader.start", "notify.Watch", err)
		return err
	}

	go func(loader *instanceLoader, c chan notify.EventInfo) {
		for true {
			select {
			case stop := <-loader.stopChan:
				if stop {
					fmt.Println("instanceLoader Stopped")
					notify.Stop(c)
					loader.stopChan <- true
				}
			case ei := <-c:
				switch ei.Event() {
				case notify.Create:
					go func() {
						var i uint
						conf := config.GetEventServiceConfig().PluginConfig
						for i = 0; i < conf.NewInstanceTryTimes; i++ {
							err = loader.newInstance(ei.Path())
							if err != nil {
								time.Sleep(time.Duration(conf.NewInstanceTimeoutSec) * time.Second)
								continue
							}

							break
						}

						if i == conf.NewInstanceTryTimes {
							fmt.Println(ei.Path(), "create failed:", err.Error())
						} else {
							fmt.Println(ei.Path(), "created")
						}
					}()
				case notify.Remove:
					fallthrough
				case notify.Rename:
					err := loader.destroyInstance(ei.Path())
					if err != nil {
						utils.PrintCallErr("instanceLoader.start", "loader.destroyInstance", err)
						continue
					}

					fmt.Println(ei.Path() + " destroyed")
				}
			}
		}
	}(loader, c)

	return nil
}

func (loader *instanceLoader) stop() {
	loader.stopChan <- true
	<-loader.stopChan
}

func (loader *instanceLoader) newInstance(instanceFilePath string) error {
	if utils.IsStringEmpty(instanceFilePath) {
		utils.PrintErr("newInstance", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	instance := loader.getInstance(instanceFilePath)
	if instance != nil {
		utils.PrintErr("newInstance", "该名称的实例已存在")
		return errors.New("该名称的实例已存在")
	}

	instanceFileName := filepath.Base(instanceFilePath)
	instanceName := strings.ReplaceAll(instanceFileName, filepath.Ext(instanceFileName), "")
	instance, err := loader.plugin.NewInstance(instanceName)
	if err != nil {
		utils.PrintCallErr("newInstance", "loader.plugin.NewInstance", err)
		return err
	}

	err = instance.Start()
	if err != nil {
		utils.PrintCallErr("newInstance", "instance.Start", err)
		return err
	}

	loader.addInstance(instanceFilePath, instance)

	return nil
}

func (loader *instanceLoader) destroyInstance(instanceFilePath string) error {
	if utils.IsStringEmpty(instanceFilePath) {
		utils.PrintErr("destroyInstance", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	instance := loader.getInstance(instanceFilePath)
	if instance == nil {
		utils.PrintErr("destroyInstance", "没有找到对应的生产者")
		return errors.New("没有找到对应的生产者")
	}

	loader.deleteInstance(instanceFilePath)

	err := instance.Stop()
	if err != nil {
		utils.PrintCallErr("destroyInstance", "instance.Stop", err)
		return err
	}

	err = loader.plugin.DestroyInstance(instance)
	if err != nil {
		utils.PrintCallErr("destroyInstance", "loader.plugin.DestroyInstance", err)
		return err
	}

	return nil
}

func (loader *instanceLoader) addInstance(instanceFilePath string, instance Instance) {
	loader.instanceMapLock.Lock()
	defer loader.instanceMapLock.Unlock()

	loader.instanceMap[instanceFilePath] = instance
}

func (loader *instanceLoader) getInstance(instanceFilePath string) Instance {
	loader.instanceMapLock.Lock()
	defer loader.instanceMapLock.Unlock()

	return loader.instanceMap[instanceFilePath]
}

func (loader *instanceLoader) deleteInstance(instanceFilePath string) {
	loader.instanceMapLock.Lock()
	defer loader.instanceMapLock.Unlock()

	loader.instanceMap[instanceFilePath] = nil
	delete(loader.instanceMap, instanceFilePath)
}

func (loader *instanceLoader) getAllInstances() []string {
	loader.instanceMapLock.Lock()
	defer loader.instanceMapLock.Unlock()

	instanceNames := make([]string, 0)
	for instanceName, _ := range loader.instanceMap {
		instanceNames = append(instanceNames, instanceName)
	}

	return instanceNames
}
