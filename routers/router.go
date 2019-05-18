package routers

import (
	"shop/pkg/setting"
	v1 "shop/routers/api/v1"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "shop/docs"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gin.SetMode(setting.RunMode)

	apia1 := r.Group("/api/v1")
	{
		//基本信息
		apia1.GET("/app", v1.GetAppBase)
		apia1.GET("/index", v1.GetAppInfo)
		apia1.GET("/detail", v1.GetGoodDetail)
		apia1.GET("/category", v1.GetGoodCategory)
		apia1.GET("/list", v1.GetGoodList)
	}

	return r
}
