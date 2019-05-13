package v1

import (
	"github.com/gin-gonic/gin"
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
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

	data["wxapp"] = models.GetAppBase(app.WxappID)
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
