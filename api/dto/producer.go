package dto

import "com.fs/event-service/service/model"

type AddProducerRequest struct {
	PluginName   string `json:"pluginName" binding:"required"`
	ProducerName string `json:"producerName" binding:"required"`
	Config       string `json:"config"`
}

type GetProducersResponse struct {
	MsgResponse
	Producers []ProducerInfoWithID
}

type ProducerInfo struct {
	ProducerName string `json:"producerName" binding:"required"`
}

type ProducerInfoWithID struct {
	ID uint64 `json:"id" binding:"required"`
	ProducerInfo
}

func FormProducerInfo(p *model.ProducerInfo) *ProducerInfo {
	if p == nil {
		return &ProducerInfo{}
	}

	return &ProducerInfo{
		ProducerName: p.ProducerName,
	}
}

func FormProducerInfoWithID(p *model.ProducerInfo) *ProducerInfoWithID {
	if p == nil {
		return &ProducerInfoWithID{}
	}

	return &ProducerInfoWithID{
		ID:           p.ID,
		ProducerInfo: *FormProducerInfo(p),
	}
}

func FormProducerInfoWithIDBatch(ps []model.ProducerInfo) []ProducerInfoWithID {
	producerInfos := make([]ProducerInfoWithID, 0)
	if ps == nil {
		return producerInfos
	}

	for _, p := range ps {
		producerInfos = append(producerInfos, *FormProducerInfoWithID(&p))
	}

	return producerInfos
}
