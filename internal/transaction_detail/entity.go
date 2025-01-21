package transactiondetail

import "time"

type TransactionDetail struct {
	Id            int       `json:"id"`
	TransactionId int       `json:"transaction_id"`
	ProductId     int       `json:"product_id"`
	Price         float64   `json:"price"`
	Qty           int       `json:"qty"`
	Subtotal      float64   `json:"subtotal"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
