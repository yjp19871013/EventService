package router

import (
	"com.fs/event-service/api"
	"github.com/gin-gonic/gin"
)

var (
	v2PostRouter = map[string]gin.HandlerFunc{
		"/load/producer-plugin": api.LoadPlugin,
	}

	v2DeleteRouter = map[string]gin.HandlerFunc{
		"/unload/producer-plugin/:id": api.UnloadPlugin,
	}

	v2PutRouter = map[string]gin.HandlerFunc{}

	v2GetRouter = map[string]gin.HandlerFunc{}
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
