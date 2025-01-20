package notification

import (
	"net/http"
	"post-backend/internal/helper"

	"github.com/gin-gonic/gin"
)

type NotificationHandler interface {
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	UpdateStatus(c *gin.Context)
}

type NotificationHandlerImpl struct {
	NotificationService NotificationService
}

func NewNotificationHandler(notificationService NotificationService) NotificationHandler {
	return &NotificationHandlerImpl{NotificationService: notificationService}
}

func (n *NotificationHandlerImpl) GetAll(c *gin.Context) {
	notifications, err := n.NotificationService.GetAll(c.Request.Context())
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get all data notification",
		Data:    notifications,
	})
}

func (n *NotificationHandlerImpl) GetById(c *gin.Context) {
	var input GetByIdNotificationInput
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	notification, err := n.NotificationService.GetById(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get all data notification",
		Data:    notification,
	})
}

func (n *NotificationHandlerImpl) UpdateStatus(c *gin.Context) {
	var inputId GetByIdNotificationInput
	if !helper.BindAndValidate(c, &inputId, "uri") {
		return
	}

	var inputData UpdateNotificationInput
	if !helper.BindAndValidate(c, &inputData, "json") {
		return
	}

	err := n.NotificationService.UpdateStatus(c.Request.Context(), inputId, inputData)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get all data notification",
		Data:    "OK",
	})
}
