package api

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/service"
	"com.fs/event-service/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
)

// AddProducer godoc
// @Summary 添加事件生产者
// @Description 添加事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param AddProducerRequest body dto.AddProducerRequest true "添加事件生产者信息"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/add/producer [post]
func AddProducer(c *gin.Context) {
	request := &dto.AddProducerRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
		return
	}

	err = service.NewProducer(request.PluginName, request.ProducerName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.AddProducer(request.PluginName, request.ProducerName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "添加事件生产者成功")
}

// DeleteProducer godoc
// @Summary 删除事件生产者
// @Description 删除事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param producerName path string true "事件生产者名称"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/delete/producer/{producerName} [delete]
func DeleteProducer(c *gin.Context) {
	producerName := c.Param("producerName")
	if utils.IsStringEmpty(producerName) {
		dto.Response400Json(c, errors.New("没有传递producerName"))
		return
	}

	err := service.DeleteProducer(producerName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DestroyProducer(producerName)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "删除事件生产者成功")
}
