package api

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/service"
	"com.fs/event-service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAllProducers godoc
// @Summary 获取事件生产者
// @Description 获取事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Success 200 {object} dto.GetProducersResponse
// @Failure 400 {object} dto.GetProducersResponse
// @Failure 500 {object} dto.GetProducersResponse
// @Router /event/api/v1/producers [get]
func GetAllProducers(c *gin.Context) {
	producers := service.GetAllProducers()
	c.JSON(http.StatusOK, dto.GetProducersResponse{
		ProducerInfos: dto.FormProducerInfoBatch(producers),
	})
}

// GetPluginAllProducers godoc
// @Summary 获取对应插件的事件生产者
// @Description 获取对应插件的事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param pluginName path string true "插件名称"
// @Success 200 {object} dto.GetPluginProducersResponse
// @Failure 400 {object} dto.GetPluginProducersResponse
// @Failure 500 {object} dto.GetPluginProducersResponse
// @Router /event/api/v1/plugin/{pluginName}/producers [get]
func GetPluginAllProducers(c *gin.Context) {
	pluginName := c.Param("pluginName")
	if utils.IsStringEmpty(pluginName) {
		c.JSON(http.StatusBadRequest, dto.GetPluginProducersResponse{
			ProducerInfo: *dto.FormProducerInfo(nil),
		})
		return
	}

	producerInfo, err := service.GetPluginProducers(pluginName)
	if err != nil {
		c.JSON(http.StatusOK, dto.GetPluginProducersResponse{
			ProducerInfo: *dto.FormProducerInfo(nil),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetPluginProducersResponse{
		ProducerInfo: *dto.FormProducerInfo(producerInfo),
	})
}
