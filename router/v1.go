package router

import (
	"com.fs/event-service/api"
	"github.com/gin-gonic/gin"
)

var (
	v1PostRouter = map[string]gin.HandlerFunc{
		"/load/producer-plugin": api.LoadPluginService,

		"/add/producer-plugin": api.AddProducerPlugin,
		"/add/producer":        api.AddProducer,
		"/add/consumer":        api.AddConsumer,
	}

	v1DeleteRouter = map[string]gin.HandlerFunc{
		"/unload/producer-plugin/:id": api.UnloadPluginService,

		"/delete/producer-plugin/:id": api.DeleteProducerPlugin,
		"/delete/producer/:id":        api.DeleteProducer,
		"/delete/consumer/:id":        api.DeleteConsumer,
	}

	v1PutRouter = map[string]gin.HandlerFunc{}

	v1GetRouter = map[string]gin.HandlerFunc{
		"/producer-plugins":                     api.GetProducerPlugins,
		"/producer-plugins/:pluginId/producers": api.GetPluginProducers,
		"/producers":                            api.GetProducers,
		"/producer/:producerId/consumers":       api.GetProducerConsumers,
		"/consumers":                            api.GetConsumers,
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
