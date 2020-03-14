package router

import (
	"com.fs/event-service/api"
	"github.com/gin-gonic/gin"
)

var (
	v2PostRouter = map[string]gin.HandlerFunc{
		"/add/producer-plugin": api.AddProducerPlugin,
		"/add/producer":        api.AddProducer,
	}

	v2DeleteRouter = map[string]gin.HandlerFunc{
		"/delete/producer-plugin/:pluginName":                        api.DeleteProducerPlugin,
		"/delete/producer-plugin/:pluginName/producer/:producerName": api.DeleteProducer,
	}

	v2PutRouter = map[string]gin.HandlerFunc{}

	v2GetRouter = map[string]gin.HandlerFunc{}
)

func initV2Router(r *gin.Engine) {
	groupV2 := r.Group(urlPrefix + "/api/v2")

	for path, f := range v2GetRouter {
		groupV2.GET(path, f)
	}

	for path, f := range v2PostRouter {
		groupV2.POST(path, f)
	}

	for path, f := range v2DeleteRouter {
		groupV2.DELETE(path, f)
	}

	for path, f := range v2PutRouter {
		groupV2.PUT(path, f)
	}
}
