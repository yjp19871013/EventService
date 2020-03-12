package http_client

import (
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	httpClient http.Client

	timeoutSec int64
}

func NewHttpClient(timeoutSec int64) *HTTPClient {
	client := new(HTTPClient)
	client.timeoutSec = timeoutSec

	if client.timeoutSec != 0 {
		client.httpClient.Timeout = time.Duration(client.timeoutSec) * time.Second
	}

	return client
}

func DestroyHttpClient(client *HTTPClient) {
	if client == nil {
		return
	}

	client.timeoutSec = 0
	client = nil
}

func (client *HTTPClient) Post(url string, request string, contentType string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(request))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return client.httpClient.Do(req)
}
