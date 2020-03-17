package producer

import (
	"com.fs/event-service/event-producer/http-push-producer/base"
	"net/http"
)

type HttpPushProducer struct {
	base.Producer
}

func HandlerFun(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
