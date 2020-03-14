package service

import (
	"com.fs/event-service/db"
	"com.fs/event-service/service/model"
	"com.fs/event-service/utils"
	"errors"
)

func AddConsumer(producerName string, consumerName string, url string) error {
	if utils.IsStringEmpty(producerName) || utils.IsStringEmpty(consumerName) {
		utils.PrintErr("AddConsumer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{Name: producerName}
	err := producer.GetByName()
	if err != nil {
		utils.PrintCallErr("AddConsumer", "p.GetByName", err)
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

func DeleteConsumer(producerName string, consumerName string) error {
	if utils.IsStringEmpty(producerName) || utils.IsStringEmpty(consumerName) {
		utils.PrintErr("DeleteConsumer", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{Name: producerName}
	err := producer.GetByName()
	if err != nil {
		utils.PrintCallErr("DeleteConsumer", "producer.GetByName", err)
		return err
	}

	consumer := &db.Consumer{Name: consumerName, ProducerID: producer.ID}
	err = consumer.GetByNameAndProducerID()
	if err != nil {
		utils.PrintCallErr("DeleteConsumer", "producer.GetByNameAndProducerID", err)
		return err
	}

	err = consumer.DeleteByIDAndName()
	if err != nil {
		utils.PrintCallErr("DeleteConsumer", "consumer.DeleteByIDAndName", err)
		return err
	}

	return nil
}

func DeleteAllConsumers(producerName string) error {
	if utils.IsStringEmpty(producerName) {
		utils.PrintErr("DeleteAllConsumers", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{Name: producerName}
	err := producer.GetByName()
	if err != nil {
		utils.PrintCallErr("DeleteAllConsumers", "producer.GetByName", err)
		return err
	}

	consumers, err := producer.GetAllProducerConsumers()
	if err != nil {
		utils.PrintCallErr("DeleteAllConsumers", "producer.GetAllProducerConsumers", err)
		return err
	}

	for _, consumer := range consumers {
		err := consumer.DeleteByIDAndName()
		if err != nil {
			utils.PrintCallErr("DeleteAllConsumers", "consumer.DeleteByIDAndName", err)
			return err
		}
	}

	return nil
}

func GetProducerConsumers(producerName string) ([]model.ConsumerInfo, error) {
	if utils.IsStringEmpty(producerName) {
		utils.PrintErr("GetProducerConsumers", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	producer := &db.Producer{Name: producerName}
	err := producer.GetByName()
	if err != nil {
		utils.PrintCallErr("GetProducerConsumers", "producer.GetByName", err)
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
