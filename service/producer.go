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
	"strconv"
)

func AddProducer(pluginID uint64, producerName string, conf string) error {
	if pluginID == 0 || utils.IsStringEmpty(producerName) {
		utils.PrintErr("AddProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{ID: pluginID}
	err := p.GetByID()
	if err != nil {
		utils.PrintCallErr("AddProducer", "p.GetByID", err)
		return err
	}

	producer := &db.Producer{
		Name:             producerName,
		Config:           conf,
		ProducerPluginID: p.ID,
	}

	err = producer.Create()
	if err != nil {
		utils.PrintCallErr("AddProducer", "producer.Create", err)
		return err
	}

	return nil
}

func DeleteProducer(producerID uint64) error {
	if producerID == 0 {
		utils.PrintErr("DeleteProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{ID: producerID}
	err := producer.GetByID()
	if err != nil {
		utils.PrintCallErr("DeleteProducer", "producer.GetByID", err)
		return err
	}

	err = producer.DeleteByID()
	if err != nil {
		utils.PrintCallErr("DeleteProducer", "producer.DeleteByID", err)
		return err
	}

	return nil
}

func DeletePluginProducers(pluginID uint64) error {
	if pluginID == 0 {
		utils.PrintErr("DeleteAllProducers", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{ID: pluginID}
	err := p.GetByID()
	if err != nil {
		utils.PrintCallErr("DeleteAllProducers", "p.GetByID", err)
		return err
	}

	producers, err := p.GetAllPluginProducers()
	if err != nil {
		utils.PrintCallErr("DeleteAllProducers", "p.GetAllPluginProducers", err)
		return err
	}

	for _, producer := range producers {
		consumers, err := producer.GetAllProducerConsumers()
		if err != nil {
			utils.PrintCallErr("DeleteAllProducers", "producer.GetAllProducerConsumers", err)
			return err
		}

		for _, consumer := range consumers {
			err := consumer.DeleteByID()
			if err != nil {
				utils.PrintCallErr("DeleteAllProducers", "producer.DeleteByIDAndName", err)
				return err
			}
		}

		err = producer.DeleteByID()
		if err != nil {
			utils.PrintCallErr("DeleteAllProducers", "p.DeleteByID", err)
			return err
		}
	}

	return nil
}

func GetPluginProducers(pluginID uint64) ([]model.ProducerInfo, error) {
	if pluginID == 0 {
		utils.PrintErr("GetPluginProducers", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	p := &db.ProducerPlugin{ID: pluginID}
	err := p.GetByID()
	if err != nil {
		utils.PrintCallErr("GetPluginProducers", "p.GetByID", err)
		return nil, err
	}

	producers, err := p.GetAllPluginProducers()
	if err != nil {
		utils.PrintCallErr("GetPluginProducers", "p.GetAllPluginProducers", err)
		return nil, err
	}

	return model.TransferProducerToProducerInfoBatch(producers), err
}

func GetAllProducers() ([]model.ProducerInfo, error) {
	producers, err := db.GetAllProducers()
	if err != nil {
		utils.PrintCallErr("GetAllProducers", "db.GetAllProducers", err)
		return nil, err
	}

	return model.TransferProducerToProducerInfoBatch(producers), nil
}

func NewProducerService(producerID uint64) error {
	if producerID == 0 {
		utils.PrintErr("NewProducerService", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	conf := config.GetEventServiceConfig().ServicesConfig
	for _, baseUrl := range conf.BaseUrls {
		url := baseUrl + "/api/v2/new/producer"
		request := dto.NewProducerRequest{ID: producerID}
		requestJson, err := json.Marshal(request)
		if err != nil {
			utils.PrintCallErr("NewProducerService", "json.Marshal", err)
			return err
		}

		client := http_client.NewHttpClient(config.HttpTimeoutSec)
		response, err := client.Post(url, string(requestJson), "application/json")
		if err != nil {
			utils.PrintCallErr("NewProducerService", "client.Post", err)
			return err
		}

		if response.StatusCode != http.StatusOK {
			utils.PrintErr("NewProducerService", baseUrl+": 响应失败")
			return errors.New(baseUrl + ": 响应失败")
		}

		responseByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			utils.PrintCallErr("NewProducerService", "ioutil.ReadAll", err)
			return err
		}

		loadPluginResponse := &dto.MsgResponse{}
		err = json.Unmarshal(responseByte, loadPluginResponse)
		if err != nil {
			utils.PrintCallErr("NewProducerService", "json.Unmarshal", err)
			return err
		}

		if !loadPluginResponse.Success {
			utils.PrintErr("NewProducerService", loadPluginResponse.Msg)
			return errors.New(loadPluginResponse.Msg)
		}
	}

	return nil
}

func NewProducer(producerID uint64) error {
	if producerID == 0 {
		utils.PrintErr("NewProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{ID: producerID}
	err := producer.GetByID()
	if err != nil {
		utils.PrintCallErr("NewProducer", "producer.GetByID", err)
		return err
	}

	p := &db.ProducerPlugin{ID: producer.ProducerPluginID}
	err = p.GetByID()
	if err != nil {
		utils.PrintCallErr("NewProducer", "p.GetByID", err)
		return err
	}

	err = loader.newProducer(p.PluginFileName, producer.Name, producer.Config)
	if err != nil {
		utils.PrintCallErr("NewProducer", "loader.newProducer", err)
		return err
	}

	return nil
}

func DestroyProducerService(producerID uint64) error {
	if producerID == 0 {
		utils.PrintErr("DestroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producerIDStr := strconv.FormatUint(producerID, 10)

	conf := config.GetEventServiceConfig().ServicesConfig
	for _, baseUrl := range conf.BaseUrls {
		url := baseUrl + "/api/v2/destroy/producer/" + producerIDStr

		client := http_client.NewHttpClient(config.HttpTimeoutSec)
		response, err := client.Delete(url, "application/json")
		if err != nil {
			utils.PrintCallErr("DestroyProducerService", "client.Delete", err)
			return err
		}

		if response.StatusCode != http.StatusOK {
			utils.PrintErr("DestroyProducerService", baseUrl+": 响应失败")
			return errors.New(baseUrl + ": 响应失败")
		}

		responseByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			utils.PrintCallErr("DestroyProducerService", "ioutil.ReadAll", err)
			return err
		}

		destroyProducerResponse := &dto.MsgResponse{}
		err = json.Unmarshal(responseByte, destroyProducerResponse)
		if err != nil {
			utils.PrintCallErr("DestroyProducerService", "json.Unmarshal", err)
			return err
		}

		if !destroyProducerResponse.Success {
			utils.PrintErr("DestroyProducerService", destroyProducerResponse.Msg)
			return errors.New(destroyProducerResponse.Msg)
		}
	}

	return nil
}

func DestroyProducer(producerID uint64) error {
	if producerID == 0 {
		utils.PrintErr("DestroyProducer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{ID: producerID}
	err := producer.GetByID()
	if err != nil {
		utils.PrintCallErr("DestroyProducer", "producer.GetByID", err)
		return err
	}

	p := &db.ProducerPlugin{ID: producer.ProducerPluginID}
	err = p.GetByID()
	if err != nil {
		utils.PrintCallErr("DestroyProducer", "p.GetByID", err)
		return err
	}

	err = loader.destroyProducer(p.PluginFileName, producer.Name)
	if err != nil {
		utils.PrintCallErr("DestroyProducer", "loader.destroyProducer", err)
		return err
	}

	return nil
}

func GetCreatedProducersService() (map[string][]string, error) {
	retMap := make(map[string][]string)

	conf := config.GetEventServiceConfig().ServicesConfig
	for _, baseUrl := range conf.BaseUrls {
		url := baseUrl + "/api/v2/created/producers"

		client := http_client.NewHttpClient(config.HttpTimeoutSec)
		response, err := client.Get(url)
		if err != nil {
			utils.PrintCallErr("GetCreatedProducersService", "client.Get", err)
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			utils.PrintErr("GetCreatedProducersService", baseUrl+": 响应失败")
			return nil, errors.New(baseUrl + ": 响应失败")
		}

		responseByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			utils.PrintCallErr("GetCreatedProducersService", "ioutil.ReadAll", err)
			return nil, err
		}

		createdProducersResponse := &dto.GetCreatedProducersResponse{}
		err = json.Unmarshal(responseByte, createdProducersResponse)
		if err != nil {
			utils.PrintCallErr("GetCreatedProducersService", "json.Unmarshal", err)
			return nil, err
		}

		if !createdProducersResponse.Success {
			utils.PrintErr("GetCreatedProducersService", createdProducersResponse.Msg)
			return nil, errors.New(createdProducersResponse.Msg)
		}

		retMap[baseUrl] = createdProducersResponse.ProducerNames
	}

	return retMap, nil
}

func GetCreatedProducers() []string {
	return loader.getAllProducers()
}
