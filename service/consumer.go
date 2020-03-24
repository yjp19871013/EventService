package service

import (
	"com.fs/event-service/db"
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"errors"
)

func AddConsumer(producerName string, consumerName string, url string) error {
	if utils.IsStringEmpty(producerName) || utils.IsStringEmpty(consumerName) || utils.IsStringEmpty(url) {
		utils.PrintErr("AddConsumer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	exist, err := loader.CheckInstanceExist(producerName)
	if err != nil {
		utils.PrintCallErr("AddConsumer", "loader.CheckInstanceExist", err)
		return err
	}

	if !exist {
		utils.PrintErr("AddConsumer", "生产者不存在")
		return errors.New("生产者不存在")
	}

	consumer := &db.Consumer{
		Name:         consumerName,
		Url:          url,
		ProducerName: producerName,
	}

	err = consumer.Create()
	if err != nil {
		utils.PrintCallErr("AddConsumer", "consumer.Create", err)
		return err
	}

	return nil
}

func DeleteConsumer(id uint64) error {
	if id == 0 {
		utils.PrintErr("DeleteConsumer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	consumer := &db.Consumer{ID: id}
	err := consumer.GetByID()
	if err != nil {
		utils.PrintCallErr("DeleteConsumer", "consumer.GetByID", err)
		return err
	}

	err = consumer.DeleteByID()
	if err != nil {
		utils.PrintCallErr("DeleteConsumer", "consumer.DeleteByIDAndName", err)
		return err
	}

	return nil
}

func GetProducerConsumers(producerName string) ([]model.ConsumerInfo, error) {
	if utils.IsStringEmpty(producerName) {
		utils.PrintErr("GetProducerConsumers", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	consumers, err := db.GetConsumersByProducerName(producerName)
	if err != nil {
		utils.PrintCallErr("GetProducerConsumers", "db.GetConsumersByProducerName", err)
		return nil, err
	}

	return model.TransferConsumerToConsumerInfoBatch(consumers), err
}

func GetAllConsumers() ([]model.ConsumerInfo, error) {
	consumers, err := db.GetAllConsumers()
	if err != nil {
		utils.PrintCallErr("GetAllConsumers", "db.GetAllConsumers", err)
		return nil, err
	}

	return model.TransferConsumerToConsumerInfoBatch(consumers), nil
}
