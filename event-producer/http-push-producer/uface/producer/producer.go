package producer

import (
	"com.fs/event-service/db"
	"com.fs/event-service/event"
	"com.fs/event-service/event-producer/http-push-producer/base"
	"com.fs/event-service/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type callbackResponse struct {
	Result  int8 `json:"result"`
	Success bool `json:"success"`
}

type HttpPushProducer struct {
	base.Producer
}

func (producer *HttpPushProducer) HandlerFun(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.PrintCallErr("uface HandlerFun", "ioutil.ReadAll", err)
		return
	}

	response := &callbackResponse{
		Result:  1,
		Success: true,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		utils.PrintCallErr("uface HandlerFun", "json.Marshal", err)
		return
	}

	_, err = w.Write(responseJson)
	if err != nil {
		utils.PrintCallErr("uface HandlerFun", "w.Write", err)
		return
	}

	dbProducer := &db.Producer{Name: producer.ProducerName}
	err = dbProducer.GetByName()
	if err != nil {
		utils.PrintCallErr("uface HandlerFun", "dbProducer.GetByName", err)
		return
	}

	consumers, err := dbProducer.GetAllProducerConsumers()
	if err != nil {
		utils.PrintCallErr("uface HandlerFun", "dbProducer.GetAllProducerConsumers", err)
		return
	}

	for _, consumer := range consumers {
		event.SendEventHttpAsync(consumer.Url, producer.ProducerName, "uface face recoganition event", string(body))
	}
}
