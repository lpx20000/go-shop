package v1

import (
	"regexp"
	"shop/pkg/e"
	"shop/pkg/util"
	"shop/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetUserAddress(c *gin.Context) {
	address := &services.Address{}
	if err := address.GetUserAddress(c.GetInt("userId")); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: address})
}

func AddAddress(c *gin.Context) {
	var (
		address    *services.Address
		err        error
		addAddress services.CommonAddress
	)

	if err = c.ShouldBindWith(&addAddress, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}

	if !(regexp.MustCompile(`^1[345678]{1}\d{9}$`).MatchString(strings.TrimSpace(addAddress.Phone))) {
		util.Response(c, util.R{Code: e.ERROR, Data: "手机号码不正确"})
		return
	}
	address = &services.Address{}
	address.UserId = c.GetInt("userId")
	address.AddAddress = addAddress
	address.AddAddress.WxappId = c.GetString("wxappId")
	if err = address.ModifyUserAddress(services.CREATE); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: "添加成功"})
}

type address struct {
	AddressId int `form:"address_id" binding:"required"  json:"address_id"`
}

func SetDefaultAddress(c *gin.Context) {
	var (
		addressInfo address
		err         error
		address     *services.Address
	)

	if err = c.ShouldBindWith(&addressInfo, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	address = &services.Address{}
	address.UserId = c.GetInt("userId")
	address.AddressId = addressInfo.AddressId

	if err = address.SetDefaultAddress(); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: "设置成功"})
}

func DeleteAddress(c *gin.Context) {
	var (
		addressInfo address
		err         error
	)

	if err = c.ShouldBindWith(&addressInfo, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	address := &services.Address{}
	address.AddressId = addressInfo.AddressId

	if err = address.DeleteAddress(); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: "删除成功"})
	return
}

func GetAddressDetail(c *gin.Context) {
	var (
		addressInfo address
		err         error
		address     *services.Address
	)

	if err = c.ShouldBindQuery(&addressInfo); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	address = &services.Address{}
	address.UserId = c.GetInt("userId")
	address.AddressId = addressInfo.AddressId

	if err = address.GetAddressDetail(); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: address.DetailAddress})
	return
}

func EditAddress(c *gin.Context) {
	var (
		address services.Address
		err     error
	)

	if err = c.ShouldBindWith(&address.UpdateAddress, binding.FormPost); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}

	if !(regexp.MustCompile(`^1[345678]{1}\d{9}$`).MatchString(strings.TrimSpace(address.UpdateAddress.Phone))) {
		util.Response(c, util.R{Code: e.ERROR, Data: "手机号码不正确"})
		return
	}
	address.UserId = c.GetInt("userId")
	address.UpdateAddress.WxappId = c.GetString("wxappId")
	if err = address.ModifyUserAddress(services.UPDATE); err != nil {
		util.Response(c, util.R{Code: e.ERROR, Data: err.Error()})
		return
	}
	util.Response(c, util.R{Code: e.SUCCESS, Data: "修改成功"})
}
