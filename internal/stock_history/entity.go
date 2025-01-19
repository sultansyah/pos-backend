package stockhistory

import "time"

type StockHistory struct {
	Id          int       `json:"id"`
	ProductId   int       `json:"product_id"`
	UserId      int       `json:"user_id"`
	Type        string    `json:"type"`
	Qty         int       `json:"qty"`
	StockBefore int       `json:"stock_before"`
	StockAfter  int       `json:"stock_after"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
