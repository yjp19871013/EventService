package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterProducer godoc
// @Summary 注册事件生产者
// @Description 注册事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param RegisterProducerRequest body dto.RegisterProducerRequest true "注册事件生产者信息"
// @Success 200 {object} dto.RegisterProducerResponse
// @Failure 400 {object} dto.RegisterProducerResponse
// @Failure 500 {object} dto.RegisterProducerResponse
// @Router /event/api/v1/register/producer [post]
func RegisterProducer(c *gin.Context) {
	c.String(http.StatusOK, "RegisterProducer OK")
}

// UnregisterProducer godoc
// @Summary 删除事件生产者
// @Description 删除事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param UnregisterProducerRequest body dto.UnregisterProducerRequest true "删除事件生产者信息"
// @Success 200 {object} dto.UnregisterProducerResponse
// @Failure 400 {object} dto.UnregisterProducerResponse
// @Failure 500 {object} dto.UnregisterProducerResponse
// @Router /event/api/v1/unregister/producer [post]
func UnregisterProducer(c *gin.Context) {
	c.String(http.StatusOK, "UnregisterProducer OK")
}
