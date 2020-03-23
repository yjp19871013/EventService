package producer

import (
	"com.fs/event-service/event-producer/http-pull-producer/base"
	"com.fs/event-service/http_client"
	"com.fs/event-service/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

	if response.StatusCode != http.StatusOK {
		utils.PrintErr("fs_relay onPull", "响应失败")
		return errors.New("响应失败")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.PrintCallErr("fs_relay onPull", "ioutil.ReadAll", err)
		return err
	}

	fmt.Println(string(body))

	//dbProducer := &db.Producer{Name: prod.ProducerName}
	//err = dbProducer.GetByName()
	//if err != nil {
	//	utils.PrintCallErr("fs_relay onPull", "dbProducer.GetByName", err)
	//	return err
	//}
	//
	//consumers, err := dbProducer.GetAllProducerConsumers()
	//if err != nil {
	//	utils.PrintCallErr("fs_relay onPull", "dbProducer.GetAllProducerConsumers", err)
	//	return err
	//}
	//
	//for _, consumer := range consumers {
	//	event.SendEventHttpAsync(consumer.Url, prod.ProducerName, "fs_relay status event", string(body))
	//}

	return nil
}
