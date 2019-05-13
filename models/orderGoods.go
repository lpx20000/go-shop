package models

type OrderGoods struct {
	OrderGoodsId    uint    `json:"order_goods_id"`
	OrderId         uint    `json:"order_id"`
	GoodsName       string  `json:"goods_name"`
	ImageId         uint    `json:"image_id"`
	DeductStockType uint8   `json:"deduct_stock_type"`
	SpecType        uint8   `json:"spec_type"`
	SpecSkuId       string  `json:"spec_sku_id"`
	GoodsAttr       string  `json:"goods_attr"`
	Content         string  `json:"content"`
	GoodsNo         string  `json:"goods_no"`
	GoodsPrice      float32 `json:"goods_price"`
	LinePrice       float32 `json:"line_price"`
	GoodsWeight     float64 `json:"goods_weight"`
	TotalNum        uint    `json:"total_num"`
	TotalPrice      float32 `json:"total_price"`
	UserId          uint    `json:"user_id"`
	WxappId         uint    `json:"wxapp_id"`
}
