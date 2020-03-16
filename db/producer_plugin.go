package db

import (
	"com.fs/event-service/utils"
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type ProducerPlugin struct {
	ID             uint64 `gorm:"primary_key;"`
	Name           string `gorm:"not null;type:varchar(256);"`
	PluginFileName string `gorm:"not null;type:varchar(256);"`

	Producers []Producer
}

func (p *ProducerPlugin) Create() error {
	if utils.IsStringEmpty(p.Name) || utils.IsStringEmpty(p.PluginFileName) {
		utils.PrintErr("ProducerPlugin.Create", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	var existCount uint64
	err := getInstance().Model(&ProducerPlugin{}).Where("name = ? OR plugin_file_name = ?",
		p.Name, p.PluginFileName).
		Count(&existCount).Error
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

func (p *ProducerPlugin) DeleteByID() error {
	if p.ID == 0 {
		utils.PrintErr("ProducerPlugin.DeleteByID", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	err := getInstance().Where("id = ?", p.ID).Delete(p).Error
	if err != nil {
		utils.PrintCallErr("DeleteByID", "Delete producer plugin", err)
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
