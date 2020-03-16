package service

import (
	"com.fs/event-service/db"
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"errors"
)

func AddConsumer(producerID uint64, consumerName string, url string) error {
	if producerID == 0 || utils.IsStringEmpty(consumerName) || utils.IsStringEmpty(url) {
		utils.PrintErr("AddConsumer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{ID: producerID}
	err := producer.GetByID()
	if err != nil {
		utils.PrintCallErr("AddConsumer", "p.GetByID", err)
		return err
	}

	consumer := &db.Consumer{
		Name:       consumerName,
		Url:        url,
		ProducerID: producer.ID,
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

func DeleteProducerConsumers(producerID uint64) error {
	if producerID == 0 {
		utils.PrintErr("DeleteAllConsumers", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{ID: producerID}
	err := producer.GetByID()
	if err != nil {
		utils.PrintCallErr("DeleteAllConsumers", "producer.GetByID", err)
		return err
	}

	consumers, err := producer.GetAllProducerConsumers()
	if err != nil {
		utils.PrintCallErr("DeleteAllConsumers", "producer.GetAllProducerConsumers", err)
		return err
	}

	for _, consumer := range consumers {
		err := consumer.DeleteByID()
		if err != nil {
			utils.PrintCallErr("DeleteAllConsumers", "consumer.DeleteByIDAndName", err)
			return err
		}
	}

	return nil
}

func GetProducerConsumers(producerID uint64) ([]model.ConsumerInfo, error) {
	if producerID == 0 {
		utils.PrintErr("GetProducerConsumers", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{ID: producerID}
	err := producer.GetByID()
	if err != nil {
		utils.PrintCallErr("GetProducerConsumers", "producer.GetByID", err)
		return nil, err
	}

	consumers, err := producer.GetAllProducerConsumers()
	if err != nil {
		utils.PrintCallErr("GetProducerConsumers", "producer.GetAllProducerConsumers", err)
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
