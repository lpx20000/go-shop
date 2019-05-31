package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin"
)

type Good struct {
	GoodId int `form:"goods_id" binding:"required" json:"goods_id"`
}

type List struct {
	Page       int    `form:"page"`
	SortType   string `form:"sortType" binding:"required"`
	SortPrice  int    `form:"sortPrice"`
	CategoryId int    `form:"category_id" binding:"required"`
	Search     string `form:"search"`
}

// @Summary 获取小程序商品详情
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Param good_id query string true "GoodId"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/detail?wxapp_id={id} [get]
func GetGoodDetail(c *gin.Context) {
	var (
		good       Good
		err        error
		goodDetail *services.GoodDetail
	)
	if err = c.ShouldBindQuery(&good); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	goodDetail = &services.GoodDetail{}
	if err = goodDetail.GetGoodDetail(good.GoodId); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: goodDetail})
}

// @Summary 获取小程序商品分类
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Param page query string true "Page"
// @Param sortType query string true "SortType"
// @Param sortPrice query string true "SortPrice"
// @Param category_id query string true "CategoryId"
// @Param search query string false "Search"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/list?wxapp_id={id} [get]
func GetGoodList(c *gin.Context) {
	var (
		list     List
		pageInfo services.GoodsListPage
	)

	if err := c.ShouldBindQuery(&list); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err})
		return
	}

	if list.Page == 0 {
		list.Page = 1
	}

	models.Host = c.Request.Host
	pageInfo = services.GoodsListPage{
		Page:       list.Page,
		CategoryId: list.CategoryId,
		SortPrice:  list.SortPrice,
		Search:     list.Search,
		SortType:   list.SortType,
	}

	if pageInfo.Page > 1 {
		pageInfo.Page = list.Page
	}

	h := &services.GoodsList{}
	if err := h.GetGoodsList(pageInfo); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: h})
}
