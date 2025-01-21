package transaction

import (
	"context"
	"database/sql"
	"post-backend/internal/custom"
	"post-backend/internal/helper"
	"post-backend/internal/payment"
	transactiondetail "post-backend/internal/transaction_detail"
)

type TransactionService interface {
	GetAll(ctx context.Context) ([]Transaction, error)
	GetById(ctx context.Context, input GetTransactionInput) (Transaction, error)
	Insert(ctx context.Context, input InsertTransactionInput) (Transaction, error)
	// Update(ctx context.Context, inputId GetTransactionInput, inputData InsertTransactionInput) error
	Delete(ctx context.Context, input GetTransactionInput) error
}

type TransactionServiceImpl struct {
	DB                          *sql.DB
	TransactionRepository       TransactionRepository
	TransactionDetailRepository transactiondetail.TransactionDetailRepository
	PaymentService              payment.PaymentService
}

func NewTransactionService(DB *sql.DB, transactionRepository TransactionRepository, transactionDetailRepository transactiondetail.TransactionDetailRepository, paymentService payment.PaymentService) TransactionService {
	return &TransactionServiceImpl{
		DB:                          DB,
		TransactionRepository:       transactionRepository,
		TransactionDetailRepository: transactionDetailRepository,
		PaymentService:              paymentService,
	}
}

func (t *TransactionServiceImpl) GetAll(ctx context.Context) ([]Transaction, error) {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return []Transaction{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	transactions, err := t.TransactionRepository.FindAll(ctx, tx)
	if err != nil {
		return []Transaction{}, err
	}

	return transactions, nil
}

func (t *TransactionServiceImpl) GetById(ctx context.Context, input GetTransactionInput) (Transaction, error) {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return Transaction{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	transaction, err := t.TransactionRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return Transaction{}, err
	}
	if transaction.Id < 0 {
		return Transaction{}, custom.ErrNotFound
	}

	return transaction, nil
}

func (t *TransactionServiceImpl) Delete(ctx context.Context, input GetTransactionInput) error {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	transaction, err := t.TransactionRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if transaction.Id < 0 {
		return custom.ErrNotFound
	}

	err = t.TransactionRepository.Delete(ctx, tx, input.Id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionServiceImpl) Insert(ctx context.Context, input InsertTransactionInput) (Transaction, error) {
	panic("unimplemented")
}
