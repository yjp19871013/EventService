package model

type ProducerInfo struct {
	PluginName    string
	ProducerNames []string
}

func TransferInstancesToProducerInfo(producerName string, instances []string) *ProducerInfo {
	if instances == nil {
		return &ProducerInfo{
			PluginName:    "",
			ProducerNames: make([]string, 0),
		}
	}

	return &ProducerInfo{
		PluginName:    producerName,
		ProducerNames: instances,
	}
}

func TransferInstancesToProducerInfoBatch(instances map[string][]string) []ProducerInfo {
	producerInfos := make([]ProducerInfo, 0)
	if instances == nil {
		return producerInfos
	}

	for pluginName, instanceNames := range instances {
		producerInfos = append(producerInfos, *TransferInstancesToProducerInfo(pluginName, instanceNames))
	}

	return producerInfos
}
