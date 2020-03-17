package router

import (
	"com.fs/event-service/api"
	"github.com/gin-gonic/gin"
)

var (
	v2PostRouter = map[string]gin.HandlerFunc{
		"/load/producer-plugin": api.LoadPlugin,
		"/new/producer":         api.NewProducer,
	}

	v2DeleteRouter = map[string]gin.HandlerFunc{
		"/unload/producer-plugin/:id": api.UnloadPlugin,
		"/destroy/producer/:id":       api.DestroyProducer,
	}

	v2PutRouter = map[string]gin.HandlerFunc{}

	v2GetRouter = map[string]gin.HandlerFunc{
		"/loaded/producer-plugins": api.GetLoadedPlugins,
		"/created/producers":       api.GetCreatedProducers,
	}
)

func initV2Router(r *gin.Engine) {
	groupV1 := r.Group(urlPrefix + "/api/v2")

	for path, f := range v2GetRouter {
		groupV1.GET(path, f)
	}

	for path, f := range v2PostRouter {
		groupV1.POST(path, f)
	}

	for path, f := range v2DeleteRouter {
		groupV1.DELETE(path, f)
	}

	for path, f := range v2PutRouter {
		groupV1.PUT(path, f)
	}
}
