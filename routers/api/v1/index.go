package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin"
)

type AppRequest struct {
	WxappID uint `form:"wxapp_id" binding:"required"`
}

// @Summary 获取小程序首页信息
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/index?wxapp_id={id} [get]
func GetAppInfo(c *gin.Context) {
	data := make(map[string]interface{})
	models.SetInfo(c.Request.Host)
	data["items"] = services.GetPageItem()
	data["newest"] = services.GetNewestGood()
	data["best"] = services.GetBestGoods()
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}

// @Summary 获取小程序基本信息
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/app?wxapp_id={id} [get]
func GetAppBase(c *gin.Context) {
	var app AppRequest
	data := make(map[string]interface{})

	if c.ShouldBindQuery(&app) != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: data})
		return
	}
	var err error
	if info, err := services.GetAppBase(app.WxappID); err == nil {
		data["wxapp"] = info
		util.Response(c, util.R{Code: e.SUCCESS, Data: data})
		return
	}
	data["wxapp"] = err
	util.Response(c, util.R{Code: e.ERROR, Data: data})
}
