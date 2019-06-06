package models

import (
	"html"

	"github.com/jinzhu/gorm"
)

const (
	ON_SALES    = 10
	SINGLE_SPEC = 10
	PER_PAGE    = 15
)

type Goods struct {
	GoodsId          int                    `gorm:"AUTO_INCREMENT,PRIMARY_KEY" json:"goods_id"`
	GoodsName        string                 `json:"goods_name"`
	CategoryId       uint                   `json:"category_id"`
	GoodsSales       uint                   `json:"goods_sales"`
	SpecType         int                    `json:"spec_type"`
	DeductStockType  int                    `json:"deduct_stock_type"`
	Content          string                 `json:"content"`
	SalesInitial     uint                   `json:"-"`
	SalesActual      uint                   `json:"-"`
	GoodsSort        uint                   `json:"goods_sort"`
	DeliveryId       uint                   `json:"delivery_id"`
	GoodsStatus      uint8                  `json:"-"`
	GoodsStatusArray map[string]interface{} `json:"goods_status"`
	IsDelete         uint8                  `json:"-"`
	WxappId          uint                   `json:"-"`
	//GoodsSkuId       string                 `json:"goods_sku_id,omitempty"`
	GoodsPrice       float64        `json:"goods_price,omitempty"`
	TotalNum         int            `json:"total_num,omitempty"`
	TotalPrice       float64        `json:"total_price,omitempty"`
	GoodsTotalWeight float64        `json:"goods_total_weight,omitempty"`
	ExpressPrice     float64        `json:"express_price,omitempty"`
	GoodsSku         GoodsSpec      `json:"goods_sku,omitempty"`
	GoodsMinPrice    string         `json:"goods_min_price,omitempty"`
	GoodsMaxPrice    string         `json:"goods_max_price,omitempty"`
	Category         Category       `gorm:"foreignkey:CategoryId;association_foreignkey:CategoryId" json:"category,omitempty" ` //belongsTo
	GoodsSpec        []GoodsSpec    `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"spec,omitempty" `           //hasMany
	GoodsImage       []GoodsImage   `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"image,omitempty" `          //hasMany
	GoodsSpecRel     []GoodsSpecRel `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"-" `                        //belongsToMany
	Delivery         Delivery       `gorm:"foreignkey:DeliveryId;association_foreignkey:DeliveryId" json:"delivery,omitempty" ` //belongsTo
	SpecRel          []SpecRel      `json:"spec_rel,omitempty"`
}

func (g *Goods) AfterFind() error {
	goodsStatus := map[uint8]map[string]interface{}{
		10: {"text": "上架", "value": 10},
		20: {"text": "下架", "value": 20},
	}
	g.GoodsStatusArray = goodsStatus[g.GoodsStatus]
	g.GoodsSales = g.SalesInitial + g.SalesActual
	g.Content = html.UnescapeString(g.Content)

	return nil
}

func GetGoodsInfoForCartList(goodsId []int) (goods []Goods) {
	Db.Where(&Goods{IsDelete: 0}).Where("goods_id in (?)", goodsId).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Preload("Delivery").
		Preload("GoodsSpecRel").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		Find(&goods)
	return
}

func GetGoodDetail(goodId int) (good Goods, err error) {
	err = Db.Where(map[string]interface{}{
		"is_delete": 0, "goods_status": ON_SALES,
		"goods_id": goodId,
	}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Preload("Delivery").
		Preload("GoodsSpecRel").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		First(&good).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

func GetIndexBestGoods() ([]*Goods, error) {
	var (
		goods []*Goods
		err   error
	)
	err = Db.Where(map[string]interface{}{"is_delete": 0, "goods_status": ON_SALES}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		Limit(10).
		Find(&goods).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return goods, err
}

func GetIndexNewestGood() ([]*Goods, error) {
	var (
		goods []*Goods
		err   error
	)
	err = Db.Where(map[string]interface{}{"is_delete": 0, "goods_status": ON_SALES}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		First(&goods).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return goods, err
}
