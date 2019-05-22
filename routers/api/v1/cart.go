package v1

import (
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
		auth   authToken
		token  interface{}
		exists bool
	)

	if c.ShouldBindQuery(&auth) != nil {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	data := make(map[string]interface{})
	if token, exists = c.Get(auth.Token); exists == false {
		util.Response(c, util.R{Code: e.INVALID_PARAMS, Data: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	services.GetCartInfo(token)

	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
