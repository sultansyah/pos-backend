package transactiondetail

import (
	"context"
	"database/sql"
)

type TransactionDetailRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, transactionDetail TransactionDetail) (TransactionDetail, error)
}

type TransactionDetailRepositoryImpl struct {
}

func NewTransactionDetailRepository() TransactionDetailRepository {
	return &TransactionDetailRepositoryImpl{}
}

func (t *TransactionDetailRepositoryImpl) Insert(ctx context.Context, tx *sql.Tx, transactionDetail TransactionDetail) (TransactionDetail, error) {
	sql := "INSERT INTO transaction_details(transaction_id, product_id, price, qty, subtotal) VALUES (?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, sql, transactionDetail.TransactionId, transactionDetail.ProductId, transactionDetail.Price, transactionDetail.Qty, transactionDetail.Subtotal)
	if err != nil {
		return transactionDetail, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return transactionDetail, err
	}

	transactionDetail.Id = int(id)
	return transactionDetail, nil
}
