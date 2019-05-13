package v1

import (
	"github.com/gin-gonic/gin"
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
)

func GetAppInfo(c *gin.Context) {
	data := make(map[string]interface{})

	data["items"] = models.GetPageItem()
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
