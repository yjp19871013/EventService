package api

import (
	"com.fs/event-service/api/dto"
	"com.fs/event-service/service"
	"com.fs/event-service/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"net/http"
	"strconv"
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

	err = service.AddProducer(request.PluginID, request.ProducerName, request.Config)
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
// @Param id path string true "事件生产者ID"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/delete/producer/{id} [delete]
func DeleteProducer(c *gin.Context) {
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

	err = service.DeleteProducerConsumers(id)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DeleteProducer(id)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "删除事件生产者成功")
}

// GetPluginProducers godoc
// @Summary 获取某个插件下的所有事件生产者
// @Description 获取某个插件下的所有事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param pluginId path string true "插件ID"
// @Success 200 {object} dto.GetProducersResponse
// @Failure 400 {object} dto.GetProducersResponse
// @Failure 500 {object} dto.GetProducersResponse
// @Router /event/api/v1/producer-plugins/{pluginId}/producers [get]
func GetPluginProducers(c *gin.Context) {
	pluginIDStr := c.Param("pluginId")
	if utils.IsStringEmpty(pluginIDStr) {
		c.JSON(http.StatusBadRequest, dto.GetProducersResponse{
			MsgResponse: dto.FormFailureMsgResponse("获取某个插件下的所有事件生产者失败", errors.New("没有传递pluginId")),
			Producers:   dto.FormProducerInfoWithIDBatch(nil),
		})
		return
	}

	pluginID, err := strconv.ParseUint(pluginIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.GetProducersResponse{
			MsgResponse: dto.FormFailureMsgResponse("获取某个插件下的所有事件生产者失败", err),
			Producers:   dto.FormProducerInfoWithIDBatch(nil),
		})
		return
	}

	producers, err := service.GetPluginProducers(pluginID)
	if err != nil {
		c.JSON(http.StatusOK, dto.GetProducersResponse{
			MsgResponse: dto.FormFailureMsgResponse("获取某个插件下的所有事件生产者失败", err),
			Producers:   dto.FormProducerInfoWithIDBatch(nil),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetProducersResponse{
		MsgResponse: dto.FormSuccessMsgResponse("获取某个插件下的所有事件生产者成功"),
		Producers:   dto.FormProducerInfoWithIDBatch(producers),
	})
}

// GetProducers godoc
// @Summary 获取所有事件生产者
// @Description 获取所有事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Success 200 {object} dto.GetProducersResponse
// @Failure 400 {object} dto.GetProducersResponse
// @Failure 500 {object} dto.GetProducersResponse
// @Router /event/api/v1/producers [get]
func GetProducers(c *gin.Context) {
	producers, err := service.GetAllProducers()
	if err != nil {
		c.JSON(http.StatusOK, dto.GetProducersResponse{
			MsgResponse: dto.FormFailureMsgResponse("获取所有事件生产者失败", err),
			Producers:   dto.FormProducerInfoWithIDBatch(nil),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetProducersResponse{
		MsgResponse: dto.FormSuccessMsgResponse("获取所有事件生产者成功"),
		Producers:   dto.FormProducerInfoWithIDBatch(producers),
	})
}

// NewProducerService godoc
// @Summary 创建所有服务的生产者
// @Description 创建所有服务的生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param NewProducerRequest body dto.NewProducerRequest true "创建所有服务的生产者请求"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/new/producer [post]
func NewProducerService(c *gin.Context) {
	request := &dto.NewProducerRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
		return
	}

	err = service.NewProducerService(request.ID)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "创建生产者成功")
}

// NewProducer godoc
// @Summary 创建生产者
// @Description 创建生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param NewProducerRequest body dto.NewProducerRequest true "创建生产者请求"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v2/new/producer [post]
func NewProducer(c *gin.Context) {
	request := &dto.NewProducerRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		dto.Response400Json(c, err)
		return
	}

	err = service.NewProducer(request.ID)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "创建生产者成功")
}

// DestroyProducerService godoc
// @Summary 销毁所有服务的事件生产者
// @Description 销毁所有服务的事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param id path string true "生产者ID"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v1/destroy/producer/{id} [delete]
func DestroyProducerService(c *gin.Context) {
	idStr := c.Param("id")
	if utils.IsStringEmpty(idStr) {
		dto.Response400Json(c, errors.New("没有传递ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DestroyProducerService(id)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "销毁所有服务的事件生产者成功")
}

// DestroyProducer godoc
// @Summary 销毁事件生产者
// @Description 销毁事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Param id path string true "生产者ID"
// @Success 200 {object} dto.MsgResponse
// @Failure 400 {object} dto.MsgResponse
// @Failure 500 {object} dto.MsgResponse
// @Router /event/api/v2/destroy/producer/{id} [delete]
func DestroyProducer(c *gin.Context) {
	idStr := c.Param("id")
	if utils.IsStringEmpty(idStr) {
		dto.Response400Json(c, errors.New("没有传递ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	err = service.DestroyProducer(id)
	if err != nil {
		dto.Response200FailJson(c, err)
		return
	}

	dto.Response200Json(c, "销毁事件生产者成功")
}

// GetCreatedProducersService godoc
// @Summary 获取所有服务已创建的事件生产者
// @Description 获取所有服务已创建的事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Success 200 {object} dto.GetCreatedProducersServiceResponse
// @Failure 400 {object} dto.GetCreatedProducersServiceResponse
// @Failure 500 {object} dto.GetCreatedProducersServiceResponse
// @Router /event/api/v1/created/producers [get]
func GetCreatedProducersService(c *gin.Context) {
	createdProducersMap, err := service.GetCreatedProducersService()
	if err != nil {
		c.JSON(http.StatusOK, dto.GetCreatedProducersServiceResponse{
			MsgResponse:      dto.FormFailureMsgResponse("获取所有服务已创建的事件生产者失败", err),
			ServiceProducers: make([]dto.ServiceProducers, 0),
		})
		return
	}

	ret := make([]dto.ServiceProducers, 0)
	for baseUrl, producerNames := range createdProducersMap {
		serviceProducers := dto.ServiceProducers{
			BaseUrl:       baseUrl,
			ProducerNames: producerNames,
		}

		ret = append(ret, serviceProducers)
	}

	c.JSON(http.StatusOK, dto.GetCreatedProducersServiceResponse{
		MsgResponse:      dto.FormSuccessMsgResponse("获取所有服务已创建的事件生产者成功"),
		ServiceProducers: ret,
	})
}

// GetCreatedProducers godoc
// @Summary 获取已创建的事件生产者
// @Description 获取已创建的事件生产者
// @Tags 事件生产者
// @Accept  json
// @Produce json
// @Success 200 {object} dto.GetCreatedProducersResponse
// @Failure 400 {object} dto.GetCreatedProducersResponse
// @Failure 500 {object} dto.GetCreatedProducersResponse
// @Router /event/api/v2/created/producers [get]
func GetCreatedProducers(c *gin.Context) {
	c.JSON(http.StatusOK, dto.GetCreatedProducersResponse{
		MsgResponse:   dto.FormSuccessMsgResponse("获取已创建的事件生产者成功"),
		ProducerNames: service.GetCreatedProducers(),
	})
}
