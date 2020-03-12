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
	statusType    = "status"
	alarmType     = "alarm"
	infoType      = "info"
	customizeType = "customize"
)

const (
	sendTimeoutSec = 10
)

type event struct {
	Type     string `json:"type"`
	Producer string `json:"producer"`
	Name     string `json:"name"`
	Time     string `json:"time"`
	Data     string `json:"data"`
}

func SendStatusEventHttp(url string, producer string, name string, data string) error {
	return sendEventHttp(url, statusType, producer, name, data)
}

func SendAlarmEventHttp(url string, producer string, name string, data string) error {
	return sendEventHttp(url, alarmType, producer, name, data)
}

func SendInfoEventHttp(url string, producer string, name string, data string) error {
	return sendEventHttp(url, infoType, producer, name, data)
}

func SendCustomizeEventHttp(url string, producer string, name string, data string) error {
	return sendEventHttp(url, statusType, producer, name, data)
}

func sendEventHttp(url string, t string, producer string, name string, data string) error {
	if !isTypeRight(t) {
		utils.PrintErr("SendStatusEventHTTP", "事件类型错误")
		return errors.New("事件类型错误")
	}

	if utils.IsStringEmpty(url) || utils.IsStringEmpty(producer) ||
		utils.IsStringEmpty(name) {
		utils.PrintErr("SendStatusEventHTTP", "没有传递必要的参数")
		return errors.New("没有传递必要的参数")
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	ev := &event{Type: t, Producer: producer, Name: name, Data: data, Time: now}

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

func isTypeRight(t string) bool {
	return t == statusType || t == alarmType || t == infoType || t == customizeType
}
