package services

import (
	"math"
	"shop/models"
	"strings"
)

func GetGoodDetail(goodId uint) (goods models.Goods, err error) {
	err = models.Db.Where(map[string]interface{}{
		"is_delete": 0, "goods_status": models.ON_SALES,
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

func GetGoodsSpecRel(goodId uint) (specRelAll []models.SpecRel) {
	var goodsSpecRel []models.GoodsSpecRel
	models.Db.Model(&models.GoodsSpecRel{}).
		Where(&models.GoodsSpecRel{GoodsId: goodId}).
		Preload("Spec").
		Preload("SpecValue").
		Find(&goodsSpecRel)

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
