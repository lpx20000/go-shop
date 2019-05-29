package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin"
)

type Detail struct {
	GoodId uint `form:"good_id" binding:"required"`
}

type List struct {
	Page       int    `form:"page"`
	SortType   string `form:"sortType" binding:"required"`
	SortPrice  int    `form:"sortPrice"`
	CategoryId uint   `form:"category_id" binding:"required"`
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
	var detail Detail
	data := make(map[string]interface{})

	if c.ShouldBindQuery(&detail) != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: data})
		return
	}
	var err error

	if info, err := services.GetGoodDetail(detail.GoodId); err == nil {
		util.Response(c, util.R{Code: e.SUCCESS, Data: info})
		return
	}
	data["detail"] = err
	util.Response(c, util.R{Code: e.ERROR, Data: data})
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
		list List
		err  error
		page int
	)
	data := make(map[string]interface{})
	if err := c.ShouldBindQuery(&list); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err})
		return
	}

	models.Host = c.Request.Host

	page = 1

	if list.Page > 1 {
		page = list.Page
	}

	info, err := services.GetGoodsList(page,
		list.CategoryId, list.Search, list.SortType, list.SortPrice)
	if err == nil {
		data["list"] = info
		util.Response(c, util.R{Code: e.SUCCESS, Data: data})
		return
	}
	data["list"] = err.Error()
	util.Response(c, util.R{Code: e.ERROR, Data: data})
}
