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
	apia1.Use(middleware.Auth())
	{
		//基本信息
		apia1.GET("/app", v1.GetAppBase)
		apia1.GET("/help", v1.GetAppHelp)
		apia1.GET("/index", v1.GetAppInfo)
		apia1.GET("/detail", v1.GetGoodDetail)
		apia1.GET("/addressList", v1.GetUserAddress)
		apia1.POST("/addAddress", v1.AddAddress)
		apia1.GET("/category", v1.GetGoodCategory)
		apia1.GET("/list", v1.GetGoodList)
		apia1.GET("/cart", v1.GetCartList)
		apia1.POST("/AddCart", v1.AddCart)

		//待缓存
		apia1.GET("/userDetail", v1.GetUserDetail)
		apia1.POST("/setAddress", v1.SetDefaultAddress)
		apia1.POST("/deleteAddress", v1.DeleteAddress)
		apia1.GET("/addressDetail", v1.GetAddressDetail)
		apia1.POST("/editAddress", v1.EditAddress)
		apia1.GET("/orderAllList", v1.GetOrderList)
		apia1.GET("/address", v1.GetCartAddress)
	}

	return r
}
