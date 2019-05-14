package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"

	"github.com/gin-gonic/gin"
)

type AppRequest struct {
	WxappID uint `form:"wxapp_id" binding:"required"`
}

//获取基本信息
func GetAppBase(c *gin.Context) {
	var app AppRequest
	data := make(map[string]interface{})

	if c.ShouldBindQuery(&app) != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: data})
		return
	}
	var err error
	if info, err := models.GetAppBase(app.WxappID); err == nil {
		data["wxapp"] = info
		util.Response(c, util.R{Code: e.SUCCESS, Data: data})
		return
	}
	data["wxapp"] = err
	util.Response(c, util.R{Code: e.ERROR, Data: data})
}
