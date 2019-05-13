package models

type Goods struct {
	GoodsId         uint   `json:"goods_id"`
	GoodsName       string `json:"goods_name"`
	CategoryId      uint   `json:"category_id"`
	SpecType        uint   `json:"spec_type"`
	DeductStockType uint   `json:"deduct_stock_type"`
	Content         string `json:"content"`
	SalesInitial    uint   `json:"sales_initial"`
	SalesActual     uint   `json:"sales_actual"`
	GoodsSort       uint   `json:"goods_sort"`
	DelieveryId     uint   `json:"delievery_id"`
	GoodsStatus     uint8  `json:"goods_status"`
	IsDelete        uint8  `json:"is_delete"`
	WxappId         uint   `json:"wxapp_id"`
	GoodsSpec       []GoodsSpec
	Category        Category       `gorm:"foreignkey:CategoryId"`
	GoodsSpecRel    []GoodsSpecRel `gorm:"foreignkey:GoodsId"`
}
