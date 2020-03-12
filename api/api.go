package api

import (
	"com.fs/event-service/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Version godoc
// @Summary 获取平台版本
// @Description 获取平台版本
// @Tags 获取平台版本
// @Success 200 {string} string
// @Router /event/api/version [get]
func Version(c *gin.Context) {
	c.String(http.StatusOK, "Version: "+config.GetEventServiceConfig().Version)
}
