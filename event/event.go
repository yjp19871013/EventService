package event

import (
	"com.fs/event-service/http_client"
	"com.fs/event-service/utils"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	sendTimeoutSec = 10
)

type event struct {
	Producer string `json:"producer"`
	Name     string `json:"name"`
	Time     string `json:"time"`
	Data     string `json:"data"`
}

func SendEventHttpAsync(url string, producer string, name string, data string) {
	go func() {
		err := sendEventHttp(url, producer, name, data)
		if err != nil {
			utils.PrintCallErr(producer+" SendEventHttpAsync", "sendEventHttp", err)
			return
		}
	}()
}

func sendEventHttp(url string, producer string, name string, data string) error {
	if utils.IsStringEmpty(url) || utils.IsStringEmpty(producer) ||
		utils.IsStringEmpty(name) {
		utils.PrintErr("SendStatusEventHTTP", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	ev := &event{Producer: producer, Name: name, Data: data, Time: now}

	eventJson, err := json.Marshal(ev)
	if err != nil {
		utils.PrintCallErr("SendStatusEventHTTP", "json.Marshal", err)
		return err
	}

	client := http_client.NewHttpClient(sendTimeoutSec)
	defer http_client.DestroyHttpClient(client)

	response, err := client.Post(url, string(eventJson), "application/json")
	if err != nil {
		utils.PrintCallErr("SendStatusEventHTTP", "client.Post", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		utils.PrintErr("SendStatusEventHTTP", url+" 发送事件失败")
		return errors.New(url + " 发送事件失败" + response.Status)
	}

	return nil
}
