package router

import (
	"com.fs/event-service/api"
	"github.com/gin-gonic/gin"
)

var (
	v1PostRouter = map[string]gin.HandlerFunc{
		"/register/producer":   api.RegisterProducer,
		"/unregister/producer": api.UnregisterProducer,

		"/register/consumer":   api.RegisterConsumer,
		"/unregister/consumer": api.UnregisterConsumer,
	}

	v1DeleteRouter = map[string]gin.HandlerFunc{}

	v1PutRouter = map[string]gin.HandlerFunc{}

	v1GetRouter = map[string]gin.HandlerFunc{}
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
