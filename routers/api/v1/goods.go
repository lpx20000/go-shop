package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"

	"github.com/gin-gonic/gin"
)

type Detail struct {
	GoodId uint `form:"good_id" binding:"required"`
}

type List struct {
	Page       uint   `form:"page"`
	SortType   string `form:"sortType" binding:"required"`
	SortPrice  int    `form:"sortPrice"`
	CategoryId uint   `form:"category_id" binding:"required"`
	Search     string `form:"search"`
}

func GetGoodDetail(c *gin.Context) {
	var detail Detail
	data := make(map[string]interface{})

	if c.ShouldBindQuery(&detail) != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: data})
		return
	}
	var err error

	if info, err := models.GetGoodDetail(detail.GoodId); err == nil {
		data["detail"] = info
		info.SpecRel = models.GetGoodsSpecRel(info.GoodsId)
		data["specData"] = info.GetManySpecData()
		util.Response(c, util.R{Code: e.SUCCESS, Data: data})
		return
	}
	data["detail"] = err
	util.Response(c, util.R{Code: e.ERROR, Data: data})
}

func GetGoodList(c *gin.Context) {
	var list List
	data := make(map[string]interface{})
	if err := c.ShouldBindQuery(&list); err != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: err})
		return
	}
	var err error
	models.SetInfo(c.Request.Host)
	info, err := models.GetGoodsList(list.Page,
		list.CategoryId, list.Search, list.SortType, list.SortPrice)
	if err == nil {
		data["list"] = info
		util.Response(c, util.R{Code: e.SUCCESS, Data: data})
		return
	}
	data["list"] = err
	util.Response(c, util.R{Code: e.ERROR, Data: data})
}
