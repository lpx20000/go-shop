package models

type GoodsSpec struct {
	GoodsSpecId uint    `json:"goods_spec_id"`
	GoodsId     uint    `json:"goods_id"`
	GoodsNo     string  `json:"goods_no"`
	GoodsPrice  float32 `json:"goods_price"`
	LinePrice   float32 `json:"line_price"`
	StockNum    uint    `json:"stock_num"`
	GoodsSales  uint    `json:"goods_sales"`
	GoodsWeight float64 `json:"goods_weight"`
	WxappId     uint    `json:"_"`
	SpecSkuId   string  `json:"spec_sku_id"`
}
