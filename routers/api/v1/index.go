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
	var (
		appInfo *services.App
		err     error
	)
	models.Host = c.Request.Host
	appInfo = &services.App{}
	err = appInfo.GetIndexData()
	if err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: appInfo})
}

// @Summary 获取小程序基本信息
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/app?wxapp_id={id} [get]
func GetAppBase(c *gin.Context) {
	var (
		app     AppRequest
		appBase *services.AppBase
		err     error
	)
	if c.ShouldBindQuery(&app) != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: e.GetMsg(e.ERROR_NOT_EXIST_PARAM)})
		return
	}
	appBase = new(services.AppBase)
	if err = appBase.GetAppBase(app.WxappID); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.ERROR, Data: appBase})
}

// @Summary 获取小程序帮助
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/help?token={token} [get]
func GetAppHelp(c *gin.Context) {
	h := &services.Help{}
	if err := h.GetAppHelp(); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.ERROR, Data: h})
}
