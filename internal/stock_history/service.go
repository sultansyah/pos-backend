package stockhistory

import (
	"context"
	"database/sql"
	"post-backend/internal/helper"
)

type StockHistoryService interface {
	GetAllByProduct(ctx context.Context, input GetStockHistoryByProductIdInput) ([]StockHistory, error)
	GetById(ctx context.Context, input GetStockHistoryByIdInput) (StockHistory, error)
	Create(ctx context.Context, input CreateStockHistoryInput, userId int) (StockHistory, error)
}

type StockHistoryServiceImpl struct {
	StockHistoryRepository StockHistoryRepository
	DB                     *sql.DB
}

func NewStockHistoryService(stockHistoryRepository StockHistoryRepository, DB *sql.DB) StockHistoryService {
	return &StockHistoryServiceImpl{StockHistoryRepository: stockHistoryRepository, DB: DB}
}

func (s *StockHistoryServiceImpl) Create(ctx context.Context, input CreateStockHistoryInput, userId int) (StockHistory, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return StockHistory{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	stockHistory := StockHistory{
		UserId:      userId,
		ProductId:   input.ProductId,
		Type:        input.Type,
		Qty:         input.Qty,
		StockBefore: input.StockBefore,
		StockAfter:  input.StockAfter,
	}

	stockHistory, err = s.StockHistoryRepository.Insert(ctx, tx, stockHistory)
	if err != nil {
		return StockHistory{}, err
	}

	return stockHistory, nil
}

func (s *StockHistoryServiceImpl) GetAllByProduct(ctx context.Context, input GetStockHistoryByProductIdInput) ([]StockHistory, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return []StockHistory{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	stockHistories, err := s.StockHistoryRepository.FindAllByProduct(ctx, tx, input.Id)
	if err != nil {
		return []StockHistory{}, err
	}

	return stockHistories, nil
}

func (s *StockHistoryServiceImpl) GetById(ctx context.Context, input GetStockHistoryByIdInput) (StockHistory, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return StockHistory{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	stockHistory, err := s.StockHistoryRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return StockHistory{}, err
	}

	return stockHistory, nil
}
