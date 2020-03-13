package api

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/service"
	"com.fs/event-service/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

// AddProducerPlugin godoc
// @Summary 添加事件生产者插件
// @Description 添加事件生产者插件
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param AddProducerPluginRequest body dto.AddProducerPluginRequest true "注册事件生产者信息"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/add/producer-plugin [post]
func AddProducerPlugin(c *gin.Context) {
	request := &dto.AddProducerPluginRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
		return
	}

	err = service.AddProducerPlugin(request.PluginName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	//err = service.LoadProducerPlugin(request.PluginName)
	//if err != nil {
	//	dto.Response200FailJson(c, err)
	//	return
	//}

	dto.Response200Json(c, "添加事件生产者插件成功")
}

// DeleteRoomMaster godoc
// @Summary 删除事件生产者插件
// @Description 删除事件生产者插件
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param pluginName path string true "插件名称"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/delete/producer-plugin/{pluginName} [delete]
func DeleteProducerPlugin(c *gin.Context) {
	pluginName := c.Param("pluginName")
	if utils.IsStringEmpty(pluginName) {
		dto.Response400Json(c, errors.New("没有传递生产者插件名称"))
		return
	}

	err := service.DeleteProducerPlugin(pluginName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	//service.UnloadProducerPlugin(pluginName)

	dto.Response200Json(c, "删除生产者插件成功")
}
