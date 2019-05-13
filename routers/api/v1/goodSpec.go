package v1

type GoodsSpec struct {
	GoodsSpectId uint    `json:"goods_spect_id"`
	GoodsId      uint    `json:"goods_id"`
	GoodsNo      string  `json:"goods_no"`
	GoodsPrice   float32 `json:"goods_price"`
	LinePrice    float32 `json:"line_price"`
	StockNum     int     `json:"stock_num"`
	GoodsSales   int     `json:"goods_sales"`
	GoodsWeight  float64 `json:"goods_weight"`
	WxappId      int     `json:"wxapp_id"`
	SpecSkuId    string  `json:"spec_sku_id"`
}
