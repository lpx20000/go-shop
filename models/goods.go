package models

import "html"

type Goods struct {
	GoodsId          uint                   `json:"goods_id"`
	GoodsName        string                 `json:"goods_name"`
	CategoryId       uint                   `json:"category_id"`
	GoodsSales       uint                   `json:"goods_sales"`
	SpecType         uint                   `json:"spec_type"`
	DeductStockType  uint                   `json:"deduct_stock_type"`
	Content          string                 `json:"content"`
	SalesInitial     uint                   `json:"_"`
	SalesActual      uint                   `json:"_"`
	GoodsSort        uint                   `json:"goods_sort"`
	DeliveryId       uint                   `json:"delivery_id"`
	GoodsStatus      uint8                  `json:"_"`
	GoodsStatusArray map[string]interface{} `json:"goods_status"`
	IsDelete         uint8                  `json:"_"`
	WxappId          uint                   `json:"_"`
	Category         Category               `gorm:"foreignkey:CategoryId;association_foreignkey:CategoryId" json:"category" `
	GoodsSpec        []GoodsSpec            `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"spec" `
	//GoodsSpecRel    []GoodsSpecRel `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"spec" `
	GoodsImage []GoodsImage `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"image" `
}

func GetNewestGood() (goods Goods) {
	db.Where(&Goods{IsDelete: 0, GoodsStatus: 10}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		First(&goods)
	return
}

func GetBestGoods() (goods []Goods) {
	db.Where(&Goods{IsDelete: 0, GoodsStatus: 10}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		Limit(10).
		Find(&goods)
	return
}

func (goods *Goods) AfterFind() error {
	goodsStatus := map[uint8]map[string]interface{}{
		10: {"text": "上架", "value": 10},
		20: {"text": "下架", "value": 20},
	}
	goods.GoodsStatusArray = goodsStatus[goods.GoodsStatus]
	goods.GoodsSales = goods.SalesInitial + goods.SalesActual
	goods.Content = html.UnescapeString(goods.Content)
	return nil
}
