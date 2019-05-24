package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin"
)

type authToken struct {
	Token string `form:"token" binding:"required"`
}

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
	data["all"], data["tree"] = models.GetRegionInfo()
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
