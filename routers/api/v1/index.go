package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"

	"github.com/gin-gonic/gin"
)

func GetAppInfo(c *gin.Context) {
	data := make(map[string]interface{})
	models.SetInfo(c.Request.Host)
	data["items"] = models.GetPageItem()
	data["newest"] = models.GetNewestGood()
	data["best"] = models.GetBestGoods()
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}
