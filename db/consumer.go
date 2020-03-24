package db

import (
	"com.fs/event-service/utils"
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type Consumer struct {
	ID           uint64 `gorm:"primary_key;"`
	Name         string `gorm:"not null;type:varchar(256);"`
	ProducerName string `gorm:"not null;type:varchar(256);"`
	Url          string `gorm:"not null;type:varchar(2048);"`
}

func (consumer *Consumer) Create() error {
	if utils.IsStringEmpty(consumer.Name) || utils.IsStringEmpty(consumer.ProducerName) || utils.IsStringEmpty(consumer.Url) {
		utils.PrintErr("Consumer.Create", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	var existCount uint64
	err := getInstance().Model(&Consumer{}).Where("name = ? AND producer_name = ?",
		consumer.Name, consumer.ProducerName).
		Count(&existCount).Error
	if err != nil {
		utils.PrintCallErr("Consumer.Create", "Count exist consumer", err)
		return err
	}

	if existCount != 0 {
		utils.PrintErr("Consumer.Create", "消费者已存在")
		return errors.New("消费者已存在")
	}

	err = getInstance().Create(consumer).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			utils.PrintErr("Consumer.Create", "消费者已存在")
			return errors.New("消费者已存在")
		}

		utils.PrintCallErr("Consumer.Create", "创建消费者", err)
		return err
	}

	return nil
}

func (consumer *Consumer) GetByID() error {
	if consumer.ID == 0 {
		utils.PrintErr("Consumer.GetByID", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("id = ?", consumer.ID).First(consumer).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.PrintErr("Consumer.GetByID", "消费者不存在")
			return errors.New("消费者不存在")
		}

		utils.PrintCallErr("Consumer.GetByID", "Find consumer", err)
		return err
	}

	return nil
}

func GetConsumersByProducerName(producerName string) ([]Consumer, error) {
	if utils.IsStringEmpty(producerName) {
		utils.PrintErr("Consumer.GetByProducerName", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	consumers := make([]Consumer, 0)
	err := getInstance().Where("producer_name = ?", producerName).Find(&consumers).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.PrintErr("Consumer.GetByProducerName", "消费者不存在")
			return nil, errors.New("消费者不存在")
		}

		utils.PrintCallErr("Consumer.GetByProducerName", "Find consumer", err)
		return nil, err
	}

	return consumers, nil
}

func (consumer *Consumer) DeleteByID() error {
	if consumer.ID == 0 {
		utils.PrintErr("Consumer.DeleteByID", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("id = ?", consumer.ID).Delete(consumer).Error
	if err != nil {
		utils.PrintCallErr("DeleteByID", "Delete consumer", err)
		return err
	}

	return nil
}

func GetAllConsumers() ([]Consumer, error) {
	consumers := make([]Consumer, 0)
	err := getInstance().Find(&consumers).Error
	if err != nil {
		utils.PrintCallErr("GetAllConsumers", "find consumers", err)
		return nil, err
	}

	return consumers, nil
}
