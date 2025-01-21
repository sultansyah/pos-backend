package payment

import (
	"context"
	"post-backend/internal/transaction"
	"post-backend/internal/user"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService interface {
	GetPaymentUrl(ctx context.Context, payment Payment, user user.User, transaction transaction.Transaction) (string, error)
}

type PaymentServiceImpl struct {
}

func NewPaymentService() PaymentService {
	return &PaymentServiceImpl{}
}

func (p *PaymentServiceImpl) GetPaymentUrl(ctx context.Context, payment Payment, user user.User, transaction transaction.Transaction) (string, error) {
	midtrans.ServerKey = ""
	midtrans.Environment = midtrans.Sandbox

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Username,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.Id),
			GrossAmt: int64(transaction.Total),
		},
	}

	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return snapResp.RedirectURL, nil
}
