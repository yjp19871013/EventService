package router

import (
	"com.fs/event-service/api"
	_ "com.fs/event-service/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

const (
	urlPrefix = "/event"
)

func InitRouter(r *gin.Engine) {

	// swagger
	r.GET(urlPrefix+"/api/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Version
	r.GET(urlPrefix+"/api/version", api.Version)

	// 管理后端使用的API
	initV1Router(r)
}
