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
		//info.SpecRel = models.GetGoodsSpecRefer(info.GoodsId)
		data["specData"] = info.GetManySpecData()
		util.Response(c, util.R{Code: e.SUCCESS, Data: data})
		return
	}
	data["detail"] = err
	util.Response(c, util.R{Code: e.ERROR, Data: data})
}
