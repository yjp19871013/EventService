package dto

import "com.fs/event-service/service/model"

type AddProducerPluginRequest struct {
	ProducerPluginInfo
}

type GetProducerPluginsResponse struct {
	MsgResponse
	ProducerPlugins []ProducerPluginInfoWithID
}

type ProducerPluginInfo struct {
	PluginName string `json:"pluginName" binding:"required"`
}

type ProducerPluginInfoWithID struct {
	ID uint64 `json:"id" binding:"required"`
	ProducerPluginInfo
}

func FormProducerPluginInfo(p *model.ProducerPluginInfo) *ProducerPluginInfo {
	if p == nil {
		return &ProducerPluginInfo{}
	}

	return &ProducerPluginInfo{PluginName: p.PluginName}
}

func FormProducerPluginInfoWithID(p *model.ProducerPluginInfo) *ProducerPluginInfoWithID {
	if p == nil {
		return &ProducerPluginInfoWithID{}
	}

	return &ProducerPluginInfoWithID{
		ID:                 p.ID,
		ProducerPluginInfo: *FormProducerPluginInfo(p),
	}
}

func FormProducerPluginInfoWithIDBatch(ps []model.ProducerPluginInfo) []ProducerPluginInfoWithID {
	pluginInfos := make([]ProducerPluginInfoWithID, 0)
	if ps == nil {
		return pluginInfos
	}

	for _, p := range ps {
		pluginInfos = append(pluginInfos, *FormProducerPluginInfoWithID(&p))
	}

	return pluginInfos
}
