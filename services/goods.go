package services

import (
	"math"
	"shop/models"
	"strings"
)

func GetGoodDetail(goodId uint) (data map[string]interface{}, err error) {
	var (
		good models.Goods
	)
	data = make(map[string]interface{})
	if good, err = models.GetGoodDetail(goodId); err != nil {
		return
	}
	if good.SpecRel, err = GetGoodsSpecRel(goodId); err != nil {
		return
	}
	data["info"] = good
	data["specData"] = GetManySpecData(good)
	return
}

func GetGoodsSpecRel(goodId uint) (specRelAll []models.SpecRel, err error) {
	var goodsSpecRel []models.GoodsSpecRel
	goodsSpecRel, err = models.GetGoodSpecRel(goodId)
	if err != nil {
		return
	}

	for _, v := range goodsSpecRel {
		var spec models.SpecRel
		spec.SpecValue = v.SpecValue
		spec.Spec = v.Spec
		spec.Pivot.Id = v.Id
		spec.Pivot.GoodsId = v.GoodsId
		spec.Pivot.SpecId = v.SpecId
		spec.Pivot.SpecValueId = v.SpecValueId
		spec.Pivot.WxappId = v.WxappId
		spec.Pivot.CreateTimeStamp = v.CreateTimeStamp
		spec.Pivot.GoodsId = v.GoodsSpecRefer.GoodsId
		specRelAll = append(specRelAll, spec)
	}
	return
}

func GetGoodsList(page int,
	categoryId uint, search string, sortType string, sortPrice int) (data map[string]interface{}, err error) {

	var (
		goods []models.Goods
		total int
	)

	query := models.Db.Model(&models.Goods{}).Select(`yoshop_goods.*, (sales_initial + sales_actual) as goods_sales,MIN(goods_price) AS goods_min_price,
	MAX(goods_price) AS goods_max_price`).
		Where(map[string]interface{}{
			"is_delete": 0, "goods_status": models.ON_SALES,
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
	case "price":
		order := "goods_min_price DESC"

		if sortPrice > 0 {
			order = "goods_max_price DESC "
		}

		query = query.Order(order)
	}

	query.Count(&total)
	data = make(map[string]interface{})

	err = query.Offset(models.PER_PAGE * (page - 1)).Limit(models.PER_PAGE).Find(&goods).Error

	data["total"] = total
	data["per_page"] = models.PER_PAGE
	data["current_page"] = page
	data["last_page"] = math.Ceil(float64(total) / float64(models.PER_PAGE))
	data["data"] = goods

	return
}

func GetManySpecData(g models.Goods) (specAttrResult models.SpecAttrResult) {
	if g.SpecType == models.SINGLE_SPEC || len(g.SpecRel) == 0 || len(g.GoodsSpec) == 0 {
		return
	}
	specAttrData := make(map[uint]models.SpecAttrData)

	for _, specRel := range g.SpecRel {
		temp := specAttrData[specRel.SpecId]
		if temp.GroupId == 0 {
			temp.GroupId = specRel.Spec.SpecId
			temp.GroupName = specRel.Spec.SpecName
		}
		temp.SpecItem = append(temp.SpecItem, models.SpecItem{
			ItemId:    specRel.SpecValueId,
			SpecValue: specRel.SpecValue.SpecValue,
		})
		specAttrData[specRel.SpecId] = temp
	}

	for _, value := range specAttrData {
		specAttrResult.SpecAttr = append(specAttrResult.SpecAttr, value)
	}

	for _, list := range g.GoodsSpec {
		var specList models.SpecListData
		specList.GoodsSpecId = list.GoodsSpecId
		specList.Rows = make([]int, 0)
		specList.SpecSkuId = list.SpecSkuId
		specList.Form = models.Form{
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
