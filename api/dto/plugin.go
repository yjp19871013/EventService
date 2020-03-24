package dto

type GetLoadedPluginsResponse struct {
	MsgResponse
	Plugins []string `json:"plugins" binding:"required"`
}
