package models

type Spec struct {
	SpecId          uint   `json:"spec_id"`
	SpecName        string `json:"spec_name"`
	WxappId         uint   `json:"wxapp_id"`
	CreateTime      int64  `json:"-"`
	CreateTimeStamp string `json:"create_time"`
}

type SpecAttrResult struct {
	SpecAttr []SpecAttrData `json:"spec_attr"`
	SpecList []SpecListData `json:"spec_list"`
}

type SpecAttrData struct {
	GroupId   uint       `json:"group_id"`
	GroupName string     `json:"group_name"`
	SpecItem  []SpecItem `json:"spec_items"`
}

type SpecItem struct {
	ItemId    int    `json:"item_id"`
	SpecValue string `json:"spec_value"`
}

type SpecListData struct {
	GoodsSpecId int    `json:"goods_spec_id"`
	SpecSkuId   string `json:"spec_sku_id"`
	Rows        []int  `json:"rows"`
	Form        Form   `json:"form"`
}

type Form struct {
	GoodsNo     string  `json:"goods_no"`
	GoodsPrice  float64 `json:"goods_price"`
	GoodsWeight float64 `json:"goods_weight"`
	LinePrice   float32 `json:"line_price"`
	StockNum    int     `json:"stock_num"`
}
