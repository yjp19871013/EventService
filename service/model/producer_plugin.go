package model

import "com.fs/event-service/db"

type ProducerPluginInfo struct {
	ID         uint64
	PluginName string
}

func TransferProducerPluginToProducerPluginInfo(p *db.ProducerPlugin) *ProducerPluginInfo {
	if p == nil {
		return &ProducerPluginInfo{}
	}

	return &ProducerPluginInfo{
		ID:         p.ID,
		PluginName: p.Name,
	}
}

func TransferProducerPluginToProducerPluginInfoBatch(ps []db.ProducerPlugin) []ProducerPluginInfo {
	pluginInfos := make([]ProducerPluginInfo, 0)
	if ps == nil {
		return pluginInfos
	}

	for _, p := range ps {
		pluginInfos = append(pluginInfos, *TransferProducerPluginToProducerPluginInfo(&p))
	}

	return pluginInfos
}
