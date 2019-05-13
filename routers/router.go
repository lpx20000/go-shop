package routers

import (
	"github.com/gin-gonic/gin"
	"shop/pkg/setting"
	v1 "shop/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	apia1 := r.Group("/api/v1")
	{
		//基本信息
		apia1.GET("/app", v1.GetAppBase)
		apia1.GET("/index", v1.GetAppInfo)
	}

	return r
}
