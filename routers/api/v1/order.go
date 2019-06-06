package v1

import (
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

func GetOrderList(c *gin.Context) {
	var (
		order *services.Order
		err   error
	)

	order = &services.Order{}
	if err = c.ShouldBindQuery(&order.OrderList); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}

	models.Host = c.Request.Host
	order.UserId = c.GetInt("userId")
	if err = order.GetOrderList(order.OrderList.DataType); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: order.OrderList})
}

func GetOrderDetail(c *gin.Context) {
	var (
		order *services.Order
		err   error
	)
	order = &services.Order{}
	if err = c.ShouldBindQuery(&order.OrderDetail); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	order.UserId = c.GetInt("userId")
	models.Host = c.Request.Host
	if err = order.GetOrderDetail(); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: order.OrderDetail})
}

//取消订单
func CancelOrder(c *gin.Context) {
	var (
		order *services.Order
		err   error
	)
	order = &services.Order{}
	if err = c.ShouldBindWith(&order.OrderDetail, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	order.UserId = c.GetInt("userId")
	models.Host = c.Request.Host
	if err = order.CancelOrder(); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}

	util.Response(c, util.R{Code: e.SUCCESS, Data: "取消成功"})
}

//确认收货
func ReceiptOrder(c *gin.Context) {
	var (
		order *services.Order
		err   error
	)
	order = &services.Order{}
	if err = c.ShouldBindWith(&order.OrderDetail, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	order.UserId = c.GetInt("userId")
	if err = order.ReceiptOrder(); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}

	util.Response(c, util.R{Code: e.SUCCESS, Data: "取消成功"})
}

//立即购买
func OrderDoNow(c *gin.Context) {
	var (
		order *services.Order
		err   error
	)
	order = &services.Order{}
	if err = c.ShouldBindWith(&order.OrderBuy, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	order.UserId = c.GetInt("userId")
	order.WxappId = c.GetString("wxappId")
	order.Token = c.GetString("token")
	order.ClientIp = c.ClientIP()
	models.Host = c.Request.Host
	if err = order.DoBuyToCart(); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error(), Msg: err.Error()})
		return
	}
	if err = order.OrderBuyNow(); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error(), Msg: err.Error()})
		return
	}
	if err = order.WxPay(); err != nil {
		util.Response(c, util.R{Code: e.PAYFAIL, Data: err.Error(), Msg: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: order.OrderBuy.WxPay})
}

//购物车提交订单
func OrderCartBuy(c *gin.Context) {
	var (
		order    *services.Order
		cartList *services.UserCartList
		err      error
	)
	order = &services.Order{}
	cartList = &services.UserCartList{}

	order.UserId = c.GetInt("userId")
	order.WxappId = c.GetString("wxappId")
	order.Token = c.GetString("token")
	order.ClientIp = c.ClientIP()
	models.Host = c.Request.Host

	if err = cartList.GetCartInfo(order.WxappId, order.UserId, true); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
	}

	cartList.UserId = order.UserId
	cartList.DeleteCartCache()
	order.OrderBuy.CartOrder = cartList.CartList
	if err = order.OrderBuyNow(); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error(), Msg: err.Error()})
		return
	}
	if err = order.WxPay(); err != nil {
		util.Response(c, util.R{Code: e.PAYFAIL, Data: err.Error(), Msg: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: order.OrderBuy.WxPay})
}

//立刻提交订单
func OrderBuyNow(c *gin.Context) {
	var (
		order *services.Order
		err   error
	)
	order = &services.Order{}
	if err = c.ShouldBindQuery(&order.OrderBuy); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	order.UserId = c.GetInt("userId")
	order.WxappId = c.GetString("wxappId")
	order.Token = c.GetString("token")
	order.ClientIp = c.ClientIP()
	models.Host = c.Request.Host
	if err = order.DoBuyToCart(); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: order.OrderBuy.CartOrder})
}

func PayNotice(c *gin.Context) {

}
