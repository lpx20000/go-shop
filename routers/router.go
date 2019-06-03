package routers

import (
	"shop/middleware"
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

	gin.SetMode(setting.AppSetting.RunMode)

	apia1 := r.Group("/api/v1")
	apia1.POST("/login", v1.UserLogin)
	apia1.GET("/app", v1.GetAppBase)
	apia1.GET("/help", v1.GetAppHelp)
	apia1.GET("/index", v1.GetAppInfo)
	apia1.Use(middleware.Auth())
	{
		//基本信息
		apia1.GET("/goods/detail", v1.GetGoodDetail)
		apia1.GET("/goods/category", v1.GetGoodCategory)
		apia1.GET("/goods/list", v1.GetGoodList)
		apia1.GET("/address/list", v1.GetUserAddress)
		apia1.POST("/address/add", v1.AddAddress)
		apia1.GET("/cart/list", v1.GetCartList)
		apia1.POST("/cart/add", v1.AddCart)

		//待缓存
		apia1.POST("/cart/sub", v1.SubCart)
		apia1.POST("/cart/delete", v1.DeleteCart)
		apia1.POST("/address/set", v1.SetDefaultAddress)
		apia1.POST("/address/delete", v1.DeleteAddress)
		apia1.GET("/address/detail", v1.GetAddressDetail)
		apia1.POST("/address/edit", v1.EditAddress)
		apia1.GET("/order/list", v1.GetOrderList)
		apia1.GET("/order/detail", v1.GetOrderDetail)
		apia1.GET("/user/detail", v1.GetUserDetail)
		apia1.GET("/cart/address", v1.GetCartAddress)
	}

	return r
}
