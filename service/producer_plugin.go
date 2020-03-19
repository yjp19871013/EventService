package service

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/config"
	"com.fs/event-service/db"
	"com.fs/event-service/http_client"
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
)

func AddProducerPlugin(pluginName string, pluginFileName string) error {
	if utils.IsStringEmpty(pluginName) || utils.IsStringEmpty(pluginFileName) {
		utils.PrintErr("AddProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	conf := config.GetEventServiceConfig().PluginConfig
	pluginPath := filepath.Join(conf.Dir, pluginFileName)
	exist := utils.PathExists(pluginPath)
	if !exist {
		utils.PrintErr("AddProducerPlugin", pluginFileName+"插件文件不存在")
		return errors.New(pluginFileName + "插件文件不存在")
	}

	p := &db.ProducerPlugin{Name: pluginName, PluginFileName: pluginFileName}
	err := p.Create()
	if err != nil {
		utils.PrintCallErr("AddProducerPlugin", "p.Create", err)
		return err
	}

	return nil
}

func DeleteProducerPlugin(pluginID uint64) error {
	if pluginID == 0 {
		utils.PrintErr("DeleteProducerPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{ID: pluginID}
	err := p.GetByID()
	if err != nil {
		utils.PrintCallErr("DeleteProducerPlugin", "p.GetByID", err)
		return err
	}

	err = p.DeleteByID()
	if err != nil {
		utils.PrintCallErr("DeleteProducerPlugin", "p.DeleteByID", err)
		return err
	}

	return nil
}

func GetProducerPlugins() ([]model.ProducerPluginInfo, error) {
	plugins, err := db.GetAllProducerPlugins()
	if err != nil {
		utils.PrintCallErr("GetProducerPlugins", "db.GetAllProducerPlugins", err)
		return nil, err
	}

	return model.TransferProducerPluginToProducerPluginInfoBatch(plugins), nil
}

func LoadPluginService(pluginID uint64) error {
	if pluginID == 0 {
		utils.PrintErr("LoadPluginService", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	conf := config.GetEventServiceConfig().ServicesConfig
	for _, baseUrl := range conf.BaseUrls {
		url := baseUrl + "/api/v2/load/producer-plugin"
		request := dto.LoadPluginRequest{ID: pluginID}
		requestJson, err := json.Marshal(request)
		if err != nil {
			utils.PrintCallErr("LoadPluginService", "json.Marshal", err)
			return err
		}

		client := http_client.NewHttpClient(config.HttpTimeoutSec)
		header := map[string]string{"Content-Type": "application/json"}
		response, err := client.Post(url, string(requestJson), header)
		if err != nil {
			utils.PrintCallErr("LoadPluginService", "client.Post", err)
			return err
		}

		if response.StatusCode != http.StatusOK {
			utils.PrintErr("LoadPluginService", baseUrl+": 响应失败")
			return errors.New(baseUrl + ": 响应失败")
		}

		responseByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			utils.PrintCallErr("LoadPluginService", "ioutil.ReadAll", err)
			return err
		}

		loadPluginResponse := &dto.MsgResponse{}
		err = json.Unmarshal(responseByte, loadPluginResponse)
		if err != nil {
			utils.PrintCallErr("LoadPluginService", "json.Unmarshal", err)
			return err
		}

		if !loadPluginResponse.Success {
			utils.PrintErr("LoadPluginService", loadPluginResponse.Msg)
			return errors.New(loadPluginResponse.Msg)
		}
	}

	return nil
}

func LoadPlugin(pluginID uint64) error {
	if pluginID == 0 {
		utils.PrintErr("LoadPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{ID: pluginID}
	err := p.GetByID()
	if err != nil {
		utils.PrintCallErr("LoadPlugin", "p.GetByID", err)
		return err
	}

	err = loader.loadProducerPlugin(p.PluginFileName)
	if err != nil {
		utils.PrintCallErr("LoadPlugin", "loader.loadProducerPlugin", err)
		return err
	}

	return nil
}

func UnloadPluginService(pluginID uint64) error {
	if pluginID == 0 {
		utils.PrintErr("UnloadPluginService", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	pluginIDStr := strconv.FormatUint(pluginID, 10)

	conf := config.GetEventServiceConfig().ServicesConfig
	for _, baseUrl := range conf.BaseUrls {
		url := baseUrl + "/api/v2/unload/producer-plugin/" + pluginIDStr

		client := http_client.NewHttpClient(config.HttpTimeoutSec)
		header := map[string]string{"Content-Type": "application/json"}
		response, err := client.Delete(url, header)
		if err != nil {
			utils.PrintCallErr("UnloadPluginService", "client.Delete", err)
			return err
		}

		if response.StatusCode != http.StatusOK {
			utils.PrintErr("UnloadPluginService", baseUrl+": 响应失败")
			return errors.New(baseUrl + ": 响应失败")
		}

		responseByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			utils.PrintCallErr("UnloadPluginService", "ioutil.ReadAll", err)
			return err
		}

		unloadPluginResponse := &dto.MsgResponse{}
		err = json.Unmarshal(responseByte, unloadPluginResponse)
		if err != nil {
			utils.PrintCallErr("UnloadPluginService", "json.Unmarshal", err)
			return err
		}

		if !unloadPluginResponse.Success {
			utils.PrintErr("UnloadPluginService", unloadPluginResponse.Msg)
			return errors.New(unloadPluginResponse.Msg)
		}
	}

	return nil
}

func UnloadPlugin(pluginID uint64) error {
	if pluginID == 0 {
		utils.PrintErr("UnloadPlugin", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{ID: pluginID}
	err := p.GetByID()
	if err != nil {
		utils.PrintCallErr("UnloadPlugin", "p.GetByID", err)
		return err
	}

	loader.unloadProducerPlugin(p.PluginFileName)

	return nil
}

func GetLoadedPluginsService() (map[string][]string, error) {
	retMap := make(map[string][]string)

	conf := config.GetEventServiceConfig().ServicesConfig
	for _, baseUrl := range conf.BaseUrls {
		url := baseUrl + "/api/v2/loaded/producer-plugins"

		client := http_client.NewHttpClient(config.HttpTimeoutSec)
		response, err := client.Get(url)
		if err != nil {
			utils.PrintCallErr("GetLoadedPluginsService", "client.Get", err)
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			utils.PrintErr("GetLoadedPluginsService", baseUrl+": 响应失败")
			return nil, errors.New(baseUrl + ": 响应失败")
		}

		responseByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			utils.PrintCallErr("GetLoadedPluginsService", "ioutil.ReadAll", err)
			return nil, err
		}

		loadedPluginResponse := &dto.GetLoadedPluginsResponse{}
		err = json.Unmarshal(responseByte, loadedPluginResponse)
		if err != nil {
			utils.PrintCallErr("GetLoadedPluginsService", "json.Unmarshal", err)
			return nil, err
		}

		if !loadedPluginResponse.Success {
			utils.PrintErr("GetLoadedPluginsService", loadedPluginResponse.Msg)
			return nil, errors.New(loadedPluginResponse.Msg)
		}

		retMap[baseUrl] = loadedPluginResponse.PluginFileNames
	}

	return retMap, nil
}

func GetLoadedPlugins() []string {
	return loader.getAllProducerPlugins()
}
