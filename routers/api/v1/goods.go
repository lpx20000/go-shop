package v1

type Goods struct {
	GoodsId         uint   `json:"goods_id"`
	GoodsName       string `json:"goods_name"`
	CategoryId      uint   `json:"category_id"`
	SpecType        int    `json:"spec_type"`
	DeductStockType int    `json:"deduct_stock_type"`
	Content         string `json:"content"`
	SalesInitial    int    `json:"_"`
	SalesActual     int    `json:"_"`
	GoodsSort       int    `json:"goods_sort"`
	DeliveryId      int    `json:"delivery_id"`
	GoodsStatus     int    `json:"goods_status"`
	WxappId         uint   `json:"_"`
}
