package transaction

import (
	"context"
	"database/sql"
	"post-backend/internal/custom"
)

type TransactionRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) ([]Transaction, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (Transaction, error)
	Insert(ctx context.Context, tx *sql.Tx, transaction Transaction) (Transaction, error)
	Update(ctx context.Context, tx *sql.Tx, transaction Transaction) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type TransactionRepositoryImpl struct {
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepositoryImpl{}
}

func (t *TransactionRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql := "delete from transactions where id = ?"
	_, err := tx.ExecContext(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]Transaction, error) {
	sql := "SELECT id, user_id, total, status, code, midtrans_status, payment_status, payment_type, payment_url, note, created_at, updated_at FROM transactions"
	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return []Transaction{}, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.Id, &transaction.UserId, &transaction.Total, &transaction.Status, &transaction.Code, &transaction.MidtransStatus, &transaction.PaymentStatus, transaction.PaymentType, &transaction.PaymentURL, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return []Transaction{}, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (t *TransactionRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (Transaction, error) {
	sql := "SELECT id, user_id, total, status, code, midtrans_status, payment_status, payment_type, payment_url, note, created_at, updated_at FROM transactions"
	row, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return Transaction{}, err
	}
	defer row.Close()

	var transaction Transaction
	if row.Next() {
		if err := row.Scan(&transaction.Id, &transaction.UserId, &transaction.Total, &transaction.Status, &transaction.Code, &transaction.MidtransStatus, &transaction.PaymentStatus, transaction.PaymentType, &transaction.PaymentURL, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return Transaction{}, err
		}

		return transaction, nil
	}

	return transaction, custom.ErrNotFound
}

func (t *TransactionRepositoryImpl) Insert(ctx context.Context, tx *sql.Tx, transaction Transaction) (Transaction, error) {
	sql := "INSERT INTO transactions(user_id, total, status, code, midtrans_status, payment_status, payment_type, payment_url, note) VALUES (?,?,?,?,?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, sql, transaction.UserId, transaction.Total, transaction.Status, transaction.Code, transaction.MidtransStatus, transaction.PaymentStatus, transaction.PaymentType, transaction.PaymentURL, transaction.Note)
	if err != nil {
		return transaction, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return transaction, err
	}

	transaction.Id = int(id)
	return transaction, nil
}

func (t *TransactionRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, transaction Transaction) error {
	sql := "UPDATE transactions SET total=?,status=?,code=?,midtrans_status=?,payment_status=?,payment_type=?,payment_url=?,note=? WHERE id = ?"
	_, err := tx.ExecContext(ctx, sql, transaction.Total, transaction.Status, transaction.Code, transaction.MidtransStatus, transaction.PaymentStatus, transaction.PaymentType, transaction.PaymentURL, transaction.Id)
	if err != nil {
		return err
	}
	return nil
}
