package dto

type AddProducerRequest struct {
	PluginName   string `json:"pluginName" binding:"required"`
	ProducerName string `json:"producerName" binding:"required"`
}
