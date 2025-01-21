package transaction

type InsertTransactionInput struct {
	UserId         int     `json:"user_id" binding:"required"`
	Total          float64 `json:"total" binding:"required"`
	Status         string  `json:"status" binding:"required"` // Enum: pending, completed, canceled
	Code           string  `json:"code" binding:"required"`
	MidtransStatus string  `json:"midtrans_status" binding:"required"` // Enum: pending, settlement, expire, cancel
	PaymentStatus  string  `json:"payment_status" binding:"required"`  // Enum: unpaid, paid, failed
	PaymentType    string  `json:"payment_type" binding:"required"`    // Enum: cash, qris, bank_transfer, credit_card
	PaymentURL     string  `json:"payment_url" binding:"required"`
	Note           string  `json:"note" binding:"required"`
}

type GetTransactionInput struct {
	Id int `uri:"user_id" binding:"required"`
}
