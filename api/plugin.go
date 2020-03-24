package api

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetLoadedPlugins godoc
// @Summary 获取已加载的插件
// @Description 获取已加载的插件
// @Tags 事件插件
// @Accept  json
// @Produce json
// @Success 200 {object} dto.GetLoadedPluginsResponse
// @Failure 400 {object} dto.GetLoadedPluginsResponse
// @Failure 500 {object} dto.GetLoadedPluginsResponse
// @Router /event/api/v1/plugins [get]
func GetLoadedPlugins(c *gin.Context) {
	c.JSON(http.StatusOK, dto.GetLoadedPluginsResponse{
		MsgResponse: dto.FormSuccessMsgResponse("获取已加载的插件成功"),
		Plugins:     service.GetLoadedPlugins(),
	})
}
