package services

import (
	"encoding/json"
	"fmt"
	"math"
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
	"shop/pkg/setting"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

type GoodDetail struct {
	Detail   models.Goods           `json:"detail"`
	SpecData *models.SpecAttrResult `json:"specData"`
}

func (g *GoodDetail) getKey(goodId int) string {
	return fmt.Sprintf("%s:%d", e.CACHA_APP_GOOD, goodId)
}

func (g *GoodDetail) GetGoodDetail(goodId int) (err error) {
	var (
		good     models.Goods
		dataByte []byte
		key      string
		exist    bool
	)
	key = g.getKey(goodId)
	if dataByte, exist, err = get(key); err != nil {
		logging.LogTrace(err)
		return
	}
	if exist {
		if err = json.Unmarshal(dataByte, g); err != nil {
			logging.LogTrace(err)
		}
		return
	}
	if good, err = models.GetGoodDetail(goodId); err != nil {
		logging.LogTrace(err)
		return
	}
	if good.SpecRel, err = GetGoodsSpecRel(goodId); err != nil {
		logging.LogTrace(err)
		return
	}
	g.Detail = good
	g.SpecData = GetManySpecData(good)
	if err = set(key, g); err != nil {
		logging.LogTrace(err)
	}
	return
}

func GetGoodsSpecRel(goodId int) (specRelAll []models.SpecRel, err error) {
	var (
		goodsSpecRel []*models.GoodsSpecRel
	)
	goodsSpecRel, err = models.GetGoodSpecRel(goodId)
	if err != nil {
		return
	}

	for _, v := range goodsSpecRel {
		spec := models.SpecRel{}
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

type GoodsList struct {
	List struct {
		Total       int            `json:"total"`
		PerPage     int            `json:"per_page"`
		CurrentPage int            `json:"current_page"`
		Data        []models.Goods `json:"data"`
		LastPage    float64        `json:"last_page"`
	} `json:"list"`
}

type GoodsListPage struct {
	Page       int    `json:"page"`
	CategoryId int    `json:"category_id"`
	SortPrice  int    `json:"sort_type"`
	Search     string `json:"search"`
	SortType   string `json:"sort_type"`
}

func (g *GoodsList) GetKey(page GoodsListPage) string {
	return fmt.Sprintf("%s:%d:%d:%d:%s",
		e.CACHA_APP_GOODList, page.CategoryId, page.Page, page.SortPrice,
		strconv.Itoa(page.SortPrice)+strings.TrimSpace(page.Search))
}

func (g *GoodsList) GetGoodsList(page GoodsListPage) (err error) {
	var (
		key      string
		exist    bool
		dataByte []byte
	)
	key = g.GetKey(page)
	if dataByte, exist, err = get(key); err != nil {
		logging.LogError(err)
		return
	}
	if exist {
		err = json.Unmarshal(dataByte, g)
		return
	}
	if err = g.GetGoodsPageList(page); err != nil {
		logging.LogError(err)
		return
	}
	if len(g.List.Data) == 0 {
		return
	}
	if err = set(key, g); err != nil {
		logging.LogError(err)
	}
	return
}

func (g *GoodsList) GetGoodsPageList(page GoodsListPage) (err error) {
	var (
		goods []models.Goods
		total int
	)

	query := models.Db.Model(&models.Goods{}).Select(`yoshop_goods.*, (sales_initial + sales_actual) as goods_sales,MIN(goods_price) AS goods_min_price,
	MAX(goods_price) AS goods_max_price`).
		Where(map[string]interface{}{
			"is_delete": 0, "goods_status": models.ON_SALES,
			"category_id": page.CategoryId,
		}).
		Preload("Category").
		Preload("GoodsSpec").
		Preload("GoodsImage").
		Joins(`left join yoshop_goods_spec  on yoshop_goods.goods_id = yoshop_goods_spec.goods_id`)

	if len(strings.Trim(page.Search, "")) > 0 {
		query = query.Where("goods_name LIKE ?", "%"+page.Search+"%")
	}

	switch page.SortType {
	case "all":
		query = query.Order("goods_sort DESC").
			Order("goods_id DESC")
	case "sales":
		query = query.Order("goods_sales DESC")
	case "price":
		order := "goods_min_price DESC"

		if page.SortPrice > 0 {
			order = "goods_max_price DESC "
		}

		query = query.Order(order)
	}

	query.Count(&total)

	err = query.Offset(setting.AppSetting.PageSize * (page.Page - 1)).
		Limit(setting.AppSetting.PageSize).
		Find(&goods).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	g.List.Data = goods
	g.List.CurrentPage = page.Page
	g.List.PerPage = setting.AppSetting.PageSize
	g.List.Total = total
	g.List.LastPage = math.Ceil(float64(total) / float64(models.PER_PAGE))
	return
}

func GetManySpecData(g models.Goods) (specAttrResult *models.SpecAttrResult) {
	if g.SpecType == models.SINGLE_SPEC || len(g.SpecRel) == 0 || len(g.GoodsSpec) == 0 {
		return
	}
	specAttrData := make(map[uint]models.SpecAttrData)
	specAttrResult = &models.SpecAttrResult{}
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
