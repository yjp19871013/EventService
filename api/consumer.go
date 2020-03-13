package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddConsumer godoc
// @Summary 添加事件消费者
// @Description 添加事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Param AddConsumerRequest body dto.AddConsumerRequest true "添加事件消费者信息"
// @Success 200 {object} dto.AddConsumerResponse
// @Failure 400 {object} dto.AddConsumerResponse
// @Failure 500 {object} dto.AddConsumerResponse
// @Router /event/api/v1/add/consumer [post]
func AddConsumer(c *gin.Context) {
	c.String(http.StatusOK, "AddConsumer OK")
}

// DeleteConsumer godoc
// @Summary 删除事件消费者
// @Description 删除事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/delete/consumer/{consumerName} [delete]
func DeleteConsumer(c *gin.Context) {
	c.String(http.StatusOK, "DeleteConsumer OK")
}
