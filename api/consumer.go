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

// AddConsumer godoc
// @Summary 添加事件消费者
// @Description 添加事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Param AddProducerRequest body dto.AddConsumerRequest true "添加事件消费者信息"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/add/consumer [post]
func AddConsumer(c *gin.Context) {
	request := &dto.AddConsumerRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
		return
	}

	err = service.AddConsumer(request.ProducerName, request.ConsumerName, request.Url)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "添加事件消费者成功")
}

// DeleteConsumer godoc
// @Summary 删除事件消费者
// @Description 删除事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Param id path string true "事件生产者ID"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/delete/consumer/{id} [delete]
func DeleteConsumer(c *gin.Context) {
	idStr := c.Param("id")
	if utils.IsStringEmpty(idStr) {
		dto.Response400Json(c, errors.New("没有传递idStr"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DeleteConsumer(id)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "删除事件消费者成功")
}

// GetProducerConsumers godoc
// @Summary 获取某个生产者下的所有事件消费者
// @Description 获取某个生产者下的所有事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Param producerName path string true "事件生产者名称"
// @Success 200 {object} dto.GetConsumersResponse
// @Failure 400 {object} dto.GetConsumersResponse
// @Failure 500 {object} dto.GetConsumersResponse
// @Router /event/api/v1/producer/{producerName}/consumers [get]
func GetProducerConsumers(c *gin.Context) {
	producerName := c.Param("producerName")
	if utils.IsStringEmpty(producerName) {
		c.JSON(http.StatusBadRequest, dto.GetConsumersResponse{
			MsgResponse: dto.FormFailureMsgResponse("获取某个生产者下的所有事件消费者失败", errors.New("没有传递producerName")),
			Consumers:   dto.FormConsumerInfoWithIDBatch(nil),
		})
		return
	}

	consumers, err := service.GetProducerConsumers(producerName)
	if err != nil {
		c.JSON(http.StatusOK, dto.GetConsumersResponse{
			MsgResponse: dto.FormFailureMsgResponse("获取某个生产者下的所有事件消费者失败", err),
			Consumers:   dto.FormConsumerInfoWithIDBatch(nil),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetConsumersResponse{
		MsgResponse: dto.FormSuccessMsgResponse("获取某个生产者下的所有事件消费者成功"),
		Consumers:   dto.FormConsumerInfoWithIDBatch(consumers),
	})
}

// GetConsumers godoc
// @Summary 获取所有事件消费者
// @Description 获取所有事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Success 200 {object} dto.GetConsumersResponse
// @Failure 400 {object} dto.GetConsumersResponse
// @Failure 500 {object} dto.GetConsumersResponse
// @Router /event/api/v1/consumers [get]
func GetConsumers(c *gin.Context) {
	consumers, err := service.GetAllConsumers()
	if err != nil {
		c.JSON(http.StatusOK, dto.GetConsumersResponse{
			MsgResponse: dto.FormFailureMsgResponse("获取所有事件消费者失败", err),
			Consumers:   dto.FormConsumerInfoWithIDBatch(nil),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetConsumersResponse{
		MsgResponse: dto.FormSuccessMsgResponse("获取所有事件消费者成功"),
		Consumers:   dto.FormConsumerInfoWithIDBatch(consumers),
	})
}
