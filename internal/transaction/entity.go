package transaction

import (
	transactiondetail "post-backend/internal/transaction_detail"
	"time"
)

type Transaction struct {
	Id                 int                                   `json:"id"`
	UserId             int                                   `json:"user_id"`
	Total              float64                               `json:"total"`
	Status             string                                `json:"status"` // Enum: pending, completed, canceled
	Code               string                                `json:"code"`
	MidtransStatus     string                                `json:"midtrans_status"` // Enum: pending, settlement, expire, cancel
	PaymentStatus      string                                `json:"payment_status"`  // Enum: unpaid, paid, failed
	PaymentType        string                                `json:"payment_type"`    // Enum: cash, qris, bank_transfer, credit_card
	PaymentURL         string                                `json:"payment_url"`
	Note               string                                `json:"note"`
	CreatedAt          time.Time                             `json:"created_at"`
	UpdatedAt          time.Time                             `json:"updated_at"`
	TransactionDetails []transactiondetail.TransactionDetail `json:"transaction_details"`
}
