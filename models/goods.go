package models

import (
	"html"
)

const (
	ON_SALES    = 10
	SINGLE_SPEC = 10
	PER_PAGE    = 15
)

type Goods struct {
	GoodsId          uint                   `json:"goods_id"`
	GoodsName        string                 `json:"goods_name"`
	CategoryId       uint                   `json:"category_id"`
	GoodsSales       uint                   `json:"goods_sales"`
	SpecType         uint                   `json:"spec_type"`
	DeductStockType  uint                   `json:"deduct_stock_type"`
	Content          string                 `json:"content"`
	SalesInitial     uint                   `json:"-"`
	SalesActual      uint                   `json:"-"`
	GoodsSort        uint                   `json:"goods_sort"`
	DeliveryId       uint                   `json:"delivery_id"`
	GoodsStatus      uint8                  `json:"-"`
	GoodsStatusArray map[string]interface{} `json:"goods_status"`
	IsDelete         uint8                  `json:"-"`
	WxappId          uint                   `json:"-"`
	GoodsMinPrice    string                 `json:"goods_min_price,omitempty"`
	GoodsMaxPrice    string                 `json:"goods_max_price,omitempty"`
	Category         Category               `gorm:"foreignkey:CategoryId;association_foreignkey:CategoryId" json:"category,omitempty" ` //belongsTo
	GoodsSpec        []GoodsSpec            `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"spec,omitempty" `           //hasMany
	GoodsImage       []GoodsImage           `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"image,omitempty" `          //hasMany
	GoodsSpecRel     []GoodsSpecRel         `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"-,omitempty" `              //belongsToMany
	Delivery         Delivery               `gorm:"foreignkey:DeliveryId;association_foreignkey:DeliveryId" json:"delivery,omitempty" ` //belongsTo
	SpecRel          []SpecRel              `json:"spec_rel,omitempty"`
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

func (g *Goods) GetManySpecData() (specAttrResult SpecAttrResult) {
	if g.SpecType == SINGLE_SPEC || len(g.SpecRel) == 0 || len(g.GoodsSpec) == 0 {
		return
	}
	specAttrData := make(map[uint]SpecAttrData)

	for _, specRel := range g.SpecRel {
		temp := specAttrData[specRel.SpecId]
		if temp.GroupId == 0 {
			temp.GroupId = specRel.Spec.SpecId
			temp.GroupName = specRel.Spec.SpecName
		}
		temp.SpecItem = append(temp.SpecItem, SpecItem{
			ItemId:    specRel.SpecValueId,
			SpecValue: specRel.SpecValue.SpecValue,
		})
		specAttrData[specRel.SpecId] = temp
	}

	for _, value := range specAttrData {
		specAttrResult.SpecAttr = append(specAttrResult.SpecAttr, value)
	}

	for _, list := range g.GoodsSpec {
		var specList SpecListData
		specList.GoodsSpecId = list.GoodsSpecId
		specList.Rows = make([]int, 0)
		specList.SpecSkuId = list.SpecSkuId
		specList.Form = Form{
			GoodsNo:     list.GoodsNo,
			GoodsPrice:  list.GoodsPrice,
			GoodsWeight: list.GoodsWeight,
			LinePrice:   list.LinePrice,
			StockNum:    list.StockNum,
		}
		specAttrResult.SpecList = append(specAttrResult.SpecList, specList)
	}

	return
}
