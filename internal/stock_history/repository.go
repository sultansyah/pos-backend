package stockhistory

import (
	"context"
	"database/sql"
	"post-backend/internal/custom"
)

type StockHistoryRepository interface {
	FindAllByProduct(ctx context.Context, tx *sql.Tx, productId int) ([]StockHistory, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (StockHistory, error)
	Insert(ctx context.Context, tx *sql.Tx, stockHistory StockHistory) (StockHistory, error)
}

type StockHistoryRepositoryImpl struct {
}

func NewStockHistoryRepository() StockHistoryRepository {
	return &StockHistoryRepositoryImpl{}
}

func (s *StockHistoryRepositoryImpl) Insert(ctx context.Context, tx *sql.Tx, stockHistory StockHistory) (StockHistory, error) {
	sql := "INSERT INTO stock_history(product_id, user_id, type, qty, stock_before, stock_after) VALUES (?,?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, sql, stockHistory.ProductId, stockHistory.UserId, stockHistory.Type, stockHistory.Qty, stockHistory.StockBefore, stockHistory.StockAfter)
	if err != nil {
		return stockHistory, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return stockHistory, err
	}

	stockHistory.Id = int(id)
	return stockHistory, nil
}

func (s *StockHistoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (StockHistory, error) {
	sql := "SELECT id, product_id, user_id, type, qty, stock_before, stock_after, created_at, updated_at FROM stock_history WHERE id = ?"

	row, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return StockHistory{}, err
	}
	defer row.Close()

	var stockHistory StockHistory
	if row.Next() {
		if err := row.Scan(&stockHistory.Id, &stockHistory.ProductId, &stockHistory.UserId, &stockHistory.Type, &stockHistory.Qty, &stockHistory.StockBefore, &stockHistory.StockAfter, &stockHistory.CreatedAt, &stockHistory.UpdatedAt); err != nil {
			return StockHistory{}, err
		}

		return stockHistory, nil
	}

	return stockHistory, custom.ErrNotFound
}

func (s *StockHistoryRepositoryImpl) FindAllByProduct(ctx context.Context, tx *sql.Tx, productId int) ([]StockHistory, error) {
	sql := "SELECT id, product_id, user_id, type, qty, stock_before, stock_after, created_at, updated_at FROM stock_history WHERE product_id = ?"

	rows, err := tx.QueryContext(ctx, sql, productId)
	if err != nil {
		return []StockHistory{}, err
	}
	defer rows.Close()

	var stockHistories []StockHistory
	for rows.Next() {
		var stockHistory StockHistory
		if err := rows.Scan(&stockHistory.Id, &stockHistory.ProductId, &stockHistory.UserId, &stockHistory.Type, &stockHistory.Qty, &stockHistory.StockBefore, &stockHistory.StockAfter, &stockHistory.CreatedAt, &stockHistory.UpdatedAt); err != nil {
			return []StockHistory{}, err
		}

		stockHistories = append(stockHistories, stockHistory)
	}

	return stockHistories, nil
}
