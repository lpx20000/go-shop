package models

type GoodsSpec struct {
	GoodsSpecList
}

type GoodsSpecList struct {
	GoodsSpecId   uint    `json:"goods_spec_id"`
	GoodsId       uint    `json:"goods_id"`
	GoodsNo       string  `json:"goods_no"`
	GoodsPrice    float32 `json:"goods_price"`
	LinePrice     float32 `json:"line_price"`
	StockNum      int     `json:"stock_num"`
	GoodsSales    uint    `json:"goods_sales"`
	GoodsWeight   float64 `json:"goods_weight"`
	GoodsAttr     string  `json:"goods_attr,omitempty"`
	WxappId       uint    `json:"-"`
	SpecSkuId     string  `json:"spec_sku_id"`
	GoodsMinPrice float32 `json:"-"`
	GoodsMaxPrice float32 `json:"-"`
}

type Price struct {
	MinPrice float32 `json:"min_price"`
	MaxPrice float32 `json:"max_price"`
}
