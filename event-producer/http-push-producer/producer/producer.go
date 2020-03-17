package producer

import (
	"com.fs/event-service/utils"
	"encoding/json"
	"errors"
)

type Config struct {
	ServerUrl string `json:"serverUrl"`
	Method    string `json:"method"`
}

type Producer struct {
	config *Config
}

func InitProducer(conf string) (*Producer, error) {
	if utils.IsStringEmpty(conf) {
		utils.PrintErr("InitProducer", "没有传递配置参数")
		return nil, errors.New("没有传递配置参数")
	}

	config := &Config{}
	err := json.Unmarshal([]byte(conf), config)
	if err != nil {
		utils.PrintCallErr("InitProducer", "json.Unmarshal", err)
		return nil, err
	}

	prod := &Producer{config: config}

	return prod, nil
}

func DestroyProducer(prod *Producer) {
	if prod == nil {
		return
	}

	prod.config = nil
	prod = nil
}
