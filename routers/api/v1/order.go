package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin"
)

type orderType struct {
	Token    string `form:"token" binding:"required"`
	DataType string `form:"dataType" binding:"required"`
}

func GetOrderList(c *gin.Context) {
	var (
		token     string
		orderType orderType
		data      map[string]interface{}
		userId    int
	)
	if err := c.ShouldBindQuery(&orderType); err != nil {
		logging.LogError(err.Error())
		util.Response(c, util.R{Code: -1, Data: err.Error()})
		return
	}

	token = c.GetString(orderType.Token)
	if token == "" {
		logging.LogError(e.GetMsg(e.ERROR))
		util.Response(c, util.R{Code: e.ERROR, Data: e.GetMsg(e.ERROR)})
		return
	}
	models.Host = c.Request.Host
	userId = c.GetInt("userId")
	data = services.GetOrderList(userId, orderType.DataType)
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
