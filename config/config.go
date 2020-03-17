package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type serverConfig struct {
	Port string `json:"port"`
}

type databaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Schema   string `json:"schema"`
}

type pluginConfig struct {
	Dir               string `json:"dir"`
	ProducerConfigDir string `json:"producerConfigDir"`
}

type servicesConfig struct {
	BaseUrls []string `json:"baseUrls"`
}

type eventServiceConfig struct {
	Version        string         `json:"version"`
	ServerConfig   serverConfig   `json:"server"`
	DatabaseConfig databaseConfig `json:"database"`
	PluginConfig   pluginConfig   `json:"plugins"`
	ServicesConfig servicesConfig `json:"services"`
}

var conf = &eventServiceConfig{}

func init() {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		panic(err)
	}
}

func GetEventServiceConfig() *eventServiceConfig {
	return conf
}
