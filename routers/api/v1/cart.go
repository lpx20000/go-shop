package v1

import (
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

// @Summary 购物车列表
// @Produce  json
// @Param wxapp_id query string true "WxappID"
// @Param token query string true "Token"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success"}"
// @Success 500 {string} json "{"code":500,"data":{},"msg":"We need ID!"}"
// @Router /api/v1/list?wxapp_id={id} [get]
func GetCartList(c *gin.Context) {
	cartList := &services.UserCartList{}
	if err := cartList.GetCartInfo(c.GetString("wxappId"), c.GetInt("userId"), false); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: cartList.CartList})
}

func GetCartAddress(c *gin.Context) {
	data := make(map[string]interface{})
	//data["all"], data["tree"] = services.GetRegionInfo()
	util.Response(c, util.R{Code: e.SUCCESS, Data: data})
}

func AddCart(c *gin.Context) {
	var (
		cart services.UserCartList
		err  error
	)

	if err = c.ShouldBindWith(&cart.AddGoodCart, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}

	if err = cart.Add(c.GetInt("userId")); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: cart})
	return
}

func SubCart(c *gin.Context) {
	var (
		cart services.UserCartList
		err  error
	)

	if err = c.ShouldBindWith(&cart.SubCartList, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}

	if err = cart.Sub(c.GetInt("userId")); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: cart})
	return
}

func DeleteCart(c *gin.Context) {
	var (
		cart services.UserCartList
		err  error
	)

	if err = c.ShouldBindWith(&cart.SubCartList, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}

	if err = cart.Delete(c.GetInt("userId")); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: cart})
	return
}

//购物车结算
func OrderCart(c *gin.Context) {
	cartList := &services.UserCartList{}
	if err := cartList.GetCartInfo(c.GetString("wxappId"), c.GetInt("userId"), true); err != nil {
		util.Response(c, util.R{Code: e.FAIL, Data: err.Error()})
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: cartList.CartList})
}
