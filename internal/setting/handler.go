package setting

import (
	"net/http"
	"post-backend/internal/helper"

	"github.com/gin-gonic/gin"
)

type SettingHandler interface {
	GetAll(c *gin.Context)
}

type SettingHandlerImpl struct {
	SettingService SettingService
}

func NewSettingHandler(settingService SettingService) SettingHandler {
	return &SettingHandlerImpl{SettingService: settingService}
}

func (h *SettingHandlerImpl) GetAll(c *gin.Context) {
	settings, err := h.SettingService.GetAll(c.Request.Context())
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get all data setting",
		Data:    settings,
	})
}
