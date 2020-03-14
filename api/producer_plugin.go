package api

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/service"
	"com.fs/event-service/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddProducerPlugin godoc
// @Summary 添加事件生产者插件
// @Description 添加事件生产者插件
// @Tags 事件生产者插件
// @Accept  json
// @Produce json
// @Param AddProducerPluginRequest body dto.AddProducerPluginRequest true "添加事件生产者插件信息"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v2/add/producer-plugin [post]
func AddProducerPlugin(c *gin.Context) {
	request := &dto.AddProducerPluginRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
		return
	}

	err = service.LoadProducerPlugin(request.PluginName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.AddProducerPlugin(request.PluginName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "添加事件生产者插件成功")
}

// DeleteProducerPlugin godoc
// @Summary 删除事件生产者插件
// @Description 删除事件生产者插件
// @Tags 事件生产者插件
// @Accept  json
// @Produce json
// @Param pluginName path string true "插件名称"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v2/delete/producer-plugin/{pluginName} [delete]
func DeleteProducerPlugin(c *gin.Context) {
	pluginName := c.Param("pluginName")
	if utils.IsStringEmpty(pluginName) {
		dto.Response400Json(c, errors.New("没有传递生产者插件名称"))
		return
	}

	err := service.DeleteAllProducers(pluginName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DeleteProducerPlugin(pluginName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	service.UnloadProducerPlugin(pluginName)

	dto.Response200Json(c, "删除生产者插件成功")
}

// GetProducerPlugins godoc
// @Summary 获取所有事件生产者插件
// @Description 获取所有事件生产者插件
// @Tags 事件生产者插件
// @Accept  json
// @Produce json
// @Success 200 {object} dto.GetProducerPluginsResponse
// @Failure 400 {object} dto.GetProducerPluginsResponse
// @Failure 500 {object} dto.GetProducerPluginsResponse
// @Router /event/api/v1/producer-plugins [get]
func GetProducerPlugins(c *gin.Context) {
	ps, err := service.GetProducerPlugins()
	if err != nil {
		c.JSON(http.StatusOK, dto.GetProducerPluginsResponse{
			MsgResponse:     dto.FormFailureMsgResponse("获取所有事件生产者插件失败", err),
			ProducerPlugins: dto.FormProducerPluginInfoWithIDBatch(nil),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetProducerPluginsResponse{
		MsgResponse:     dto.FormSuccessMsgResponse("获取所有事件生产者插件成功"),
		ProducerPlugins: dto.FormProducerPluginInfoWithIDBatch(ps),
	})
}
