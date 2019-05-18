package models

import (
	"html"
	"math"
	"strings"
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
	GoodsMinPrice    float32                `json:"goods_min_price,omitempty"`
	GoodsMaxPrice    float32                `json:"goods_max_price,omitempty"`
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

func GetNewestGood() (goods Goods) {
	db.Where(map[string]interface{}{"is_delete": 0, "goods_status": ON_SALES}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Order("goods_id DESC").
		Order("goods_sort ASC").
		First(&goods)
	return
}

func GetBestGoods() (goods []Goods) {
	db.Where(map[string]interface{}{"is_delete": 0, "goods_status": ON_SALES}).
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
	err = db.Where(map[string]interface{}{
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
		First(&goods).Error
	return
}

func GetGoodsList(page int,
	categoryId uint, search string, sortType string, sortPrice int) (data map[string]interface{}, err error) {

	var (
		goods []Goods
		total int
	)

	query := db.Model(&Goods{}).Select(`yoshop_goods.*, (sales_initial + sales_actual) as goods_sales,MIN(goods_price) AS goods_min_price,
	MAX(goods_price) AS goods_max_price`).
		Where(map[string]interface{}{
			"is_delete": 0, "goods_status": ON_SALES,
			"category_id": categoryId,
		}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Joins(`left join yoshop_goods_spec  on yoshop_goods.goods_id = yoshop_goods_spec.goods_id`)

	if len(strings.Trim(search, "")) > 0 {
		query = query.Where("goods_name LIKE ?", "%"+search+"%")
	}

	switch sortType {
	case "all":
		query = query.Order("goods_sort DESC").
			Order("goods_id DESC")
	case "sales":
		query = query.Order("goods_sales DESC")
	}

	query.Count(&total)
	data = make(map[string]interface{})

	err = query.Offset(PER_PAGE * (page - 1)).Limit(PER_PAGE).Find(&goods).Error
	//if err == nil {
	//	goodsJson, _ := json.Marshal(goods)
	//	a := string(goodsJson)
	//
	//	fmt.Println(a)
	//
	//	_ = json.Unmarshal(goodsJson, &goodsList)
	//}

	data["total"] = total
	data["per_page"] = PER_PAGE
	data["current_page"] = page
	data["last_page"] = math.Ceil(float64(total) / float64(PER_PAGE))
	data["data"] = goods

	return
}
