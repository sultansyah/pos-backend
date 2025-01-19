package stockhistory

import (
	"fmt"
	"net/http"
	"post-backend/internal/helper"

	"github.com/gin-gonic/gin"
)

type StockHistoryHandler interface {
	GetAllByProduct(c *gin.Context)
	GetById(c *gin.Context)
}

type StockHistoryHandlerImpl struct {
	StockHistoryService StockHistoryService
}

func NewStockHistoryHandler(stockHistoryService StockHistoryService) StockHistoryHandler {
	return &StockHistoryHandlerImpl{StockHistoryService: stockHistoryService}
}

func (s *StockHistoryHandlerImpl) GetAllByProduct(c *gin.Context) {
	var input GetStockHistoryByProductIdInput
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	stockHistories, err := s.StockHistoryService.GetAllByProduct(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: fmt.Sprintf("success get all stock history by product id = %d", input.Id),
		Data:    stockHistories,
	})
}

func (s *StockHistoryHandlerImpl) GetById(c *gin.Context) {
	var input GetStockHistoryByIdInput
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	stockHistory, err := s.StockHistoryService.GetById(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: fmt.Sprintf("success get stock history by id = %d", input.Id),
		Data:    stockHistory,
	})
}
