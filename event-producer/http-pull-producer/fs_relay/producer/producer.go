package producer

import (
	"com.fs/event-service/db"
	"com.fs/event-service/event"
	"com.fs/event-service/event-producer/http-pull-producer/base"
	"com.fs/event-service/http_client"
	"com.fs/event-service/utils"
	"io/ioutil"
)

type HttpPullProducer struct {
	base.Producer
}

func (prod *HttpPullProducer) onPull() error {
	client := http_client.NewHttpClient(int64(prod.Config.PullTimeout))
	response, err := client.Get(prod.Config.PullUrl)
	if err != nil {
		utils.PrintCallErr("fs_relay onPull", "client.Get", err)
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.PrintCallErr("fs_relay onPull", "ioutil.ReadAll", err)
		return err
	}

	dbProducer := &db.Producer{Name: prod.ProducerName}
	err = dbProducer.GetByName()
	if err != nil {
		utils.PrintCallErr("fs_relay onPull", "dbProducer.GetByName", err)
		return err
	}

	consumers, err := dbProducer.GetAllProducerConsumers()
	if err != nil {
		utils.PrintCallErr("fs_relay onPull", "dbProducer.GetAllProducerConsumers", err)
		return err
	}

	for _, consumer := range consumers {
		event.SendEventHttpAsync(consumer.Url, prod.ProducerName, "fs_relay status event", string(body))
	}

	return nil
}
