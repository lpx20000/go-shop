package models

import (
	"html"
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
	Category         Category               `gorm:"foreignkey:CategoryId;association_foreignkey:CategoryId" json:"category" ` //belongsTo
	GoodsSpec        []GoodsSpec            `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"spec" `           //hasMany
	GoodsSpecRel     []GoodsSpecRel         `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"-" `              //belongsToMany
	GoodsImage       []GoodsImage           `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"image" `          //hasMany
	Delivery         Delivery               `gorm:"foreignkey:DeliveryId;association_foreignkey:DeliveryId" json:"delivery" ` //belongsTo
	SpecRel          []SpecRel              `json:"spec_rel"`
}

func (g *Goods) AfterFind() error {
	goodsStatus := map[uint8]map[string]interface{}{
		10: {"text": "上架", "value": 10},
		20: {"text": "下架", "value": 20},
	}
	g.GoodsStatusArray = goodsStatus[g.GoodsStatus]
	g.GoodsSales = g.SalesInitial + g.SalesActual
	g.Content = html.UnescapeString(g.Content)
	var goodsSpecRel []GoodsSpecRel
	db.Model(&GoodsSpecRel{}).
		Where(&GoodsSpecRel{GoodsId: g.GoodsId}).
		Preload("Spec").
		Preload("SpecValue").
		Find(&goodsSpecRel)

	for _, v := range goodsSpecRel {
		var spec SpecRel
		spec.SpecValue = v.SpecValue
		spec.Spec = v.Spec
		spec.Pivot.Id = v.Id
		spec.Pivot.GoodsId = v.GoodsId
		spec.Pivot.SpecId = v.SpecId
		spec.Pivot.SpecValueId = v.SpecValueId
		spec.Pivot.WxappId = v.WxappId
		spec.Pivot.CreateTimeStamp = v.CreateTimeStamp
		spec.Pivot.GoodsId = v.GoodsSpecRefer.GoodsId
		g.SpecRel = append(g.SpecRel, spec)
	}
	return nil
}

func (g *Goods) GetManySpecData() (specAttrResult SpecAttrResult) {
	if g.SpecType == 10 || len(g.SpecRel) == 0 || len(g.GoodsSpec) == 0 {
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

func GetGoodDetail(goodId uint) (goods Goods, err error) {
	err = db.Where(&Goods{IsDelete: 0, GoodsStatus: 10, GoodsId: goodId}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Preload("Delivery").
		Preload("GoodsSpecRel").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		First(&goods).Error
	return
}
