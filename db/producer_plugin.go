package db

import (
	"com.fs/event-service/utils"
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type ProducerPlugin struct {
	ID         uint64 `gorm:"primary_key;"`
	Name       string `gorm:"not null;type:varchar(256);"`
	PluginName string `gorm:"not null;type:varchar(256);"`

	Producers []Producer
}

func (p *ProducerPlugin) Create() error {
	if utils.IsStringEmpty(p.Name) || utils.IsStringEmpty(p.PluginName) {
		utils.PrintErr("ProducerPlugin.Create", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	var existCount uint64
	err := getInstance().Where("name = ? AND plugin_name = ?", p.Name, p.PluginName).Count(&existCount).Error
	if err != nil {
		utils.PrintCallErr("ProducerPlugin.Create", "Count exist plugin", err)
		return err
	}

	if existCount != 0 {
		utils.PrintErr("ProducerPlugin.Create", "生产者插件已存在")
		return errors.New("生产者插件已存在")
	}

	err = getInstance().Create(p).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			utils.PrintErr("ProducerPlugin.Create", "生产者插件已存在")
			return errors.New("生产者插件已存在")
		}

		utils.PrintCallErr("ProducerPlugin.Create", "创建生产者插件", err)
		return err
	}

	return nil
}

func (p *ProducerPlugin) GetByID() error {
	if p.ID == 0 {
		utils.PrintErr("ProducerPlugin.GetByID", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("id = ?", p.ID).First(p).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			utils.PrintErr("ProducerPlugin.GetByID", "生产者插件不存在")
			return errors.New("生产者插件不存在")
		}

		utils.PrintCallErr("ProducerPlugin.GetByID", "Find producer plugin", err)
		return err
	}

	return nil
}

func (p *ProducerPlugin) DeleteByIDAndName() error {
	if p.ID == 0 || utils.IsStringEmpty(p.Name) {
		utils.PrintErr("ProducerPlugin.DeleteByIDAndName", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("id = ? AND name = ?", p.ID, p.Name).Delete(p).Error
	if err != nil {
		utils.PrintCallErr("DeleteByIDAndName", "Delete producer plugin", err)
		return err
	}

	return nil
}

func (p *ProducerPlugin) GetAllPluginProducers() ([]Producer, error) {
	if p.ID == 0 {
		utils.PrintErr("ProducerPlugin.GetAllPluginProducers", "没有传递必要的参数")
		return nil, errors.New("没有传递必要的参数")
	}

	producers := make([]Producer, 0)
	err := getInstance().Model(p).Association("Producers").Find(&producers).Error
	if err != nil {
		utils.PrintCallErr("ProducerPlugin.GetAllPluginProducers", "find producers", err)
		return nil, err
	}

	return producers, nil
}

func GetAllProducerPlugins() ([]ProducerPlugin, error) {
	plugins := make([]ProducerPlugin, 0)
	err := getInstance().Find(&plugins).Error
	if err != nil {
		utils.PrintCallErr("GetAllProducerPlugins", "find producer plugins", err)
		return nil, err
	}

	return plugins, nil
}
