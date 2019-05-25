package v1

import (
	"github.com/gin-gonic/gin"
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"
)

func GetOrderList(c *gin.Context) {
	var (
		token  string
		auth   authToken
		data   map[string]interface{}
		userId int
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
	userId = c.GetInt("userId")
	data = services.GetOrderList(userId, "all")
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
