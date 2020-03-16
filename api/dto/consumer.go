package dto

import "com.fs/event-service/service/model"

type AddConsumerRequest struct {
	ProducerID   uint64 `json:"producerId" binding:"required"`
	ConsumerName string `json:"consumerName" binding:"required"`
	Url          string `json:"url" binding:"required"`
}

type GetConsumersResponse struct {
	MsgResponse
	Consumers []ConsumerInfoWithID
}

type ConsumerInfo struct {
	ConsumerName string `json:"consumerName" binding:"required"`
	Url          string `json:"url" binding:"required"`
}

type ConsumerInfoWithID struct {
	ID uint64 `json:"id" binding:"required"`
	ConsumerInfo
}

func FormConsumerInfo(consumer *model.ConsumerInfo) *ConsumerInfo {
	if consumer == nil {
		return &ConsumerInfo{}
	}

	return &ConsumerInfo{
		ConsumerName: consumer.Name,
		Url:          consumer.Url,
	}
}

func FormConsumerInfoWithID(consumer *model.ConsumerInfo) *ConsumerInfoWithID {
	if consumer == nil {
		return &ConsumerInfoWithID{}
	}

	return &ConsumerInfoWithID{
		ID:           consumer.ID,
		ConsumerInfo: *FormConsumerInfo(consumer),
	}
}

func FormConsumerInfoWithIDBatch(consumers []model.ConsumerInfo) []ConsumerInfoWithID {
	consumerInfos := make([]ConsumerInfoWithID, 0)
	if consumers == nil {
		return consumerInfos
	}

	for _, consumer := range consumers {
		consumerInfos = append(consumerInfos, *FormConsumerInfoWithID(&consumer))
	}

	return consumerInfos
}
