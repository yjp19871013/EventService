package router

import (
	"com.fs/event-service/api"
	"github.com/gin-gonic/gin"
)

var (
	v1PostRouter = map[string]gin.HandlerFunc{
		"/add/producer-plugin": api.AddProducerPlugin,
		"/add/producer":        api.AddProducer,
		"/add/consumer":        api.AddConsumer,
	}

	v1DeleteRouter = map[string]gin.HandlerFunc{
		"/delete/producer-plugin/:pluginName": api.DeleteProducerPlugin,
		"/delete/producer/:producerName":      api.DeleteProducer,
		"/delete/consumer/:consumerName":      api.DeleteConsumer,
	}

	v1PutRouter = map[string]gin.HandlerFunc{}

	v1GetRouter = map[string]gin.HandlerFunc{
		"/producer-plugins": api.GetProducerPlugins,
	}
)

func initV1Router(r *gin.Engine) {
	groupV1 := r.Group(urlPrefix + "/api/v1")

	for path, f := range v1GetRouter {
		groupV1.GET(path, f)
	}

	for path, f := range v1PostRouter {
		groupV1.POST(path, f)
	}

	for path, f := range v1DeleteRouter {
		groupV1.DELETE(path, f)
	}

	for path, f := range v1PutRouter {
		groupV1.PUT(path, f)
	}
}
