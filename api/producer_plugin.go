package api

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/service"
	"com.fs/event-service/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
// @Router /event/api/v1/add/producer-plugin [post]
func AddProducerPlugin(c *gin.Context) {
	request := &dto.AddProducerPluginRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
		return
	}

	err = service.AddProducerPlugin(request.PluginName, request.PluginFileName)
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
// @Param id path string true "插件ID"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/delete/producer-plugin/{id} [delete]
func DeleteProducerPlugin(c *gin.Context) {
	idStr := c.Param("id")
	if utils.IsStringEmpty(idStr) {
		dto.Response400Json(c, errors.New("没有传递id"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DeletePluginProducers(id)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DeleteProducerPlugin(id)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

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

// LoadPluginService godoc
// @Summary 加载所有服务的事件生产者插件
// @Description 加载所有服务的事件生产者插件
// @Tags 事件生产者插件
// @Accept  json
// @Produce json
// @Param LoadPluginRequest body dto.LoadPluginRequest true "加载所有服务的事件生产者插件请求"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/load/producer-plugin/{id} [post]
func LoadPluginService(c *gin.Context) {
	request := &dto.LoadPluginRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
	}

	err = service.LoadPluginService(request.ID)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "加载事件生产者插件成功")
}

// LoadPlugin godoc
// @Summary 加载事件生产者插件
// @Description 加载事件生产者插件
// @Tags 事件生产者插件
// @Accept  json
// @Produce json
// @Param LoadPluginRequest body dto.LoadPluginRequest true "加载事件生产者插件请求"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v2/load/producer-plugin/{id} [post]
func LoadPlugin(c *gin.Context) {
	request := &dto.LoadPluginRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
	}

	err = service.LoadPlugin(request.ID)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "加载事件生产者插件成功")
}
