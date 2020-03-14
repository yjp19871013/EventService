package model

import "com.fs/event-service/db"

type ConsumerInfo struct {
	ID   uint64
	Name string
	Url  string
}

func TransferConsumerToConsumerInfo(consumer *db.Consumer) *ConsumerInfo {
	if consumer == nil {
		return &ConsumerInfo{}
	}

	return &ConsumerInfo{
		ID:   consumer.ID,
		Name: consumer.Name,
		Url:  consumer.Url,
	}
}

func TransferConsumerToConsumerInfoBatch(consumers []db.Consumer) []ConsumerInfo {
	consumerInfos := make([]ConsumerInfo, 0)
	if consumers == nil {
		return consumerInfos
	}

	for _, consumer := range consumers {
		consumerInfos = append(consumerInfos, *TransferConsumerToConsumerInfo(&consumer))
	}

	return consumerInfos
}
