package models

import (
	"github.com/jinzhu/gorm"
)

type GoodsSpec struct {
	GoodsSpecId   int     `gorm:"PRIMARY_KEY" json:"goods_spec_id"`
	GoodsId       uint    `json:"goods_id"`
	GoodsNo       string  `json:"goods_no"`
	GoodsPrice    float64 `json:"goods_price"`
	LinePrice     float32 `json:"line_price"`
	StockNum      int     `json:"stock_num"`
	GoodsSales    uint    `json:"goods_sales"`
	GoodsWeight   float64 `json:"goods_weight"`
	GoodsAttr     string  `gorm:"-" json:"goods_attr,omitempty"`
	WxappId       uint    `json:"-"`
	SpecSkuId     string  `json:"spec_sku_id"`
	GoodsMinPrice float32 `gorm:"-" json:"-"`
	GoodsMaxPrice float32 `gorm:"-" json:"-"`
}

type Price struct {
	MinPrice float32 `json:"min_price"`
	MaxPrice float32 `json:"max_price"`
}

func (g *GoodsSpec) BeforeSave(scope *gorm.Scope) error {
	if g.GoodsSpecId > 0 && g.StockNum != 0 {
		stockNum := g.StockNum
		goodSpecId := g.GoodsSpecId
		_, ok := scope.FieldByName("StockNum")
		if !ok {
			return nil
		}
		err := Db.Model(&GoodsSpec{}).Where(&GoodsSpec{GoodsSpecId: goodSpecId}).First(&g).Error
		if err != nil {
			return err
		}
		g.StockNum += stockNum
		return nil
	}
	return nil
}
