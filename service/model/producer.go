package model

import "com.fs/event-service/db"

type ProducerInfo struct {
	ID           uint64
	ProducerName string
}

func TransferProducerToProducerInfo(producer *db.Producer) *ProducerInfo {
	if producer == nil {
		return &ProducerInfo{}
	}

	return &ProducerInfo{
		ID:           producer.ID,
		ProducerName: producer.Name,
	}
}

func TransferProducerToProducerInfoBatch(producers []db.Producer) []ProducerInfo {
	producerInfos := make([]ProducerInfo, 0)
	if producers == nil {
		return producerInfos
	}

	for _, p := range producers {
		producerInfos = append(producerInfos, *TransferProducerToProducerInfo(&p))
	}

	return producerInfos
}
