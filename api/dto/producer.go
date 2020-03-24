package dto

import "com.fs/event-service/service/model"

type GetProducersResponse struct {
	ProducerInfos []ProducerInfo `json:"producerInfos" binding:"required"`
}

type GetPluginProducersResponse struct {
	ProducerInfo
}

type ProducerInfo struct {
	PluginName    string
	ProducerNames []string `json:"producerNames" binding:"required"`
}

func FormProducerInfo(info *model.ProducerInfo) *ProducerInfo {
	if info == nil {
		return &ProducerInfo{
			PluginName:    "",
			ProducerNames: make([]string, 0),
		}
	}

	return &ProducerInfo{
		PluginName:    info.PluginName,
		ProducerNames: info.ProducerNames,
	}
}

func FormProducerInfoBatch(infos []model.ProducerInfo) []ProducerInfo {
	producerInfos := make([]ProducerInfo, 0)
	if infos == nil {
		return producerInfos
	}

	for _, producerInfo := range infos {
		producerInfos = append(producerInfos, *FormProducerInfo(&producerInfo))
	}

	return producerInfos
}
