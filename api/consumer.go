package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterConsumer godoc
// @Summary 注册事件消费者
// @Description 注册事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Param RegisterConsumerRequest body dto.RegisterConsumerRequest true "注册事件消费者信息"
// @Success 200 {object} dto.RegisterConsumerResponse
// @Failure 400 {object} dto.RegisterConsumerResponse
// @Failure 500 {object} dto.RegisterConsumerResponse
// @Router /event/api/v1/register/consumer [post]
func RegisterConsumer(c *gin.Context) {
	c.String(http.StatusOK, "RegisterConsumer OK")
}

// UnregisterConsumer godoc
// @Summary 删除事件消费者
// @Description 删除事件消费者
// @Tags 事件消费者
// @Accept  json
// @Produce json
// @Param UnregisterConsumerRequest body dto.UnregisterConsumerRequest true "删除事件消费者信息"
// @Success 200 {object} dto.UnregisterConsumerResponse
// @Failure 400 {object} dto.UnregisterConsumerResponse
// @Failure 500 {object} dto.UnregisterConsumerResponse
// @Router /event/api/v1/unregister/consumer [post]
func UnregisterConsumer(c *gin.Context) {
	c.String(http.StatusOK, "UnregisterConsumer OK")
}
