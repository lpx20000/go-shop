package v1

import (
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin"
)

// @Summary 购物车列表
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Param token query string true "Token"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/list?wxapp_id={id} [get]
func GetCartList(c *gin.Context) {
	var (
		auth    authToken
		token   string
		wxappId string
	)

	if c.ShouldBindQuery(&auth) != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	token = c.GetString(auth.Token)
	if token == "" {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	wxappId = c.GetString("wxappId")
	data := services.GetCartInfo(token, wxappId)

	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}

func GetCartAddress(c *gin.Context) {
	data := make(map[string]interface{})
	//data["all"], data["tree"] = services.GetRegionInfo()
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
