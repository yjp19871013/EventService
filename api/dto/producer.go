package dto

import "com.fs/event-service/service/model"

type AddProducerRequest struct {
	ProducerInfo
}

type GetProducersResponse struct {
	MsgResponse
	Producers []ProducerInfoWithID
}

type NewProducerRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

type GetCreatedProducersResponse struct {
	MsgResponse
	ProducerNames []string
}

type GetCreatedProducersServiceResponse struct {
	MsgResponse
	ServiceProducers []ServiceProducers
}

type ServiceProducers struct {
	BaseUrl       string `json:"baseUrl" binding:"required"`
	ProducerNames []string
}

type ProducerInfo struct {
	PluginID     uint64 `json:"pluginId" binding:"required"`
	ProducerName string `json:"producerName" binding:"required"`
	Config       string `json:"config"`
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
		PluginID:     p.PluginID,
		ProducerName: p.ProducerName,
		Config:       p.Config,
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
