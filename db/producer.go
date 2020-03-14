package db

import (
	"com.fs/event-service/utils"
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type Producer struct {
	ID     uint64 `gorm:"primary_key;"`
	Name   string `gorm:"not null;unique;type:varchar(256);"`
	Config string `gorm:"not null;type:varchar(10240)"`

	ProducerPluginID uint64
	Consumers        []Consumer
}

func (producer *Producer) Create() error {
	if utils.IsStringEmpty(producer.Name) || producer.ProducerPluginID == 0 {
		utils.PrintErr("Producer.Create", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Create(producer).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			utils.PrintErr("Producer.Create", "生产者已存在")
			return errors.New("生产者已存在")
		}

		utils.PrintCallErr("Producer.Create", "创建生产者", err)
		return err
	}

	return nil
}

func (producer *Producer) GetByNameAndProducerPluginID() error {
	if utils.IsStringEmpty(producer.Name) || producer.ProducerPluginID == 0 {
		utils.PrintErr("Producer.GetByNameAndProducerPluginID", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("name = ? AND producer_plugin_id = ?",
		producer.Name, producer.ProducerPluginID).
		First(producer).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.PrintErr("Producer.GetByNameAndProducerPluginID", "生产者不存在")
			return errors.New("生产者不存在")
		}

		utils.PrintCallErr("Producer.GetByNameAndProducerPluginID", "Find producer", err)
		return err
	}

	return nil
}

func (producer *Producer) DeleteByIDAndName() error {
	if producer.ID == 0 || utils.IsStringEmpty(producer.Name) {
		utils.PrintErr("Producer.DeleteByIDAndName", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("id = ? AND name = ?", producer.ID, producer.Name).Delete(producer).Error
	if err != nil {
		utils.PrintCallErr("DeleteByIDAndName", "Delete producer", err)
		return err
	}

	return nil
}

func (producer *Producer) GetAllProducerConsumers() ([]Consumer, error) {
	if producer.ID == 0 {
		utils.PrintErr("ProducerPlugin.GetAllProducerConsumers", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	consumers := make([]Consumer, 0)
	err := getInstance().Model(producer).Association("Consumers").Find(&consumers).Error
	if err != nil {
		utils.PrintCallErr("ProducerPlugin.GetAllProducerConsumers", "find consumers", err)
		return nil, err
	}

	return consumers, nil
}

func (producer *Producer) GetByName() error {
	if utils.IsStringEmpty(producer.Name) {
		utils.PrintErr("Producer.GetByName", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("name = ?", producer.Name).First(producer).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.PrintErr("Producer.GetByName", "生产者不存在")
			return errors.New("生产者不存在")
		}

		utils.PrintCallErr("Producer.GetByName", "Find producer", err)
		return err
	}

	return nil
}

func GetAllProducers() ([]Producer, error) {
	producers := make([]Producer, 0)
	err := getInstance().Find(&producers).Error
	if err != nil {
		utils.PrintCallErr("GetAllProducers", "find producers", err)
		return nil, err
	}

	return producers, nil
}
