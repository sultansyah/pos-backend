package stockhistory

type CreateStockHistoryInput struct {
	ProductId   int    `json:"product_id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Qty         int    `json:"qty" binding:"required"`
	StockBefore int    `json:"stock_before" binding:"required"`
	StockAfter  int    `json:"stock_after" binding:"required"`
}

type GetStockHistoryByIdInput struct {
	Id int `uri:"id" binding:"required"`
}

type GetStockHistoryByProductIdInput struct {
	Id int `uri:"id" binding:"required"`
}
