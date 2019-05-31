package v1

import (
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin"
)

// @Summary 获取小程序商品分类详情
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/category?wxapp_id={id} [get]
func GetGoodCategory(c *gin.Context) {
	h := &services.Category{}
	if err := h.GetCategory(); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: h})
}
