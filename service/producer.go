package service

import (
	"com.fs/event-service/db"
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"errors"
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
