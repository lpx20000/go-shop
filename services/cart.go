package services

import (
	"encoding/json"
	"fmt"
	"shop/models"
	"shop/pkg/e"
	"shop/pkg/logging"
	"shop/pkg/util"
	"strconv"
	"strings"
	"time"
)

type ErrorInfo struct {
	err    string
	exists bool
}

type UserCartList struct {
	CartList     models.CartOrder `json:"-"`
	Err          ErrorInfo        `json:"-"`
	AddGoodCart  CartList         `json:"-"`
	CartTotalNum int              `json:"cart_total_num"`
}

type CartList struct {
	GoodId     int    `form:"goods_id" binding:"required" json:"good_id"`
	GoodsNum   int    `form:"goods_num" binding:"required" json:"goods_num"`
	GoodsSkuId string `form:"goods_sku_id" binding:"required" json:"goods_sku_id"`
	CreateTime int64  `json:"create_time"`
}

func (u *UserCartList) setErrorInfo(err string) {
	u.Err.exists = true
	u.Err.err = err
}

func (u *UserCartList) getErrorInfo() (err string) {
	err = u.Err.err
	u.Err.exists = false
	u.Err.err = ""
	return
}

func (u *UserCartList) hasError() bool {
	return u.Err.exists
}

func (u *UserCartList) getKey(uid int) string {
	return e.CACHA_APP_CARTLIST + ":" + strconv.Itoa(uid)
}

func (u *UserCartList) GetCartInfo(wxappId string, uid int) (err error) {
	var (
		carts           []CartList
		goodsId         []int
		goodsInfo       map[int]models.Goods
		userInformation *models.User
		cityId          int
		existAddress    bool
		cartList        []models.Goods
		good            models.Goods
		goodSku         models.GoodsSpec
		inRegion        bool
		orderTotalPrice float64
		expressPrice    float64
		express         []float64
		OrderTotalNum   int
		key             string
		dataByte        []byte
		exist           bool
	)
	cartList = make([]models.Goods, 0)
	key = u.getKey(uid)
	if dataByte, exist, err = getDataFromRedis(key); err != nil {
		logging.LogError(err)
		return
	}

	if exist {
		if err = json.Unmarshal(dataByte, &carts); err != nil {
			logging.LogError(err)
			return
		}
	}

	for _, item := range carts {
		OrderTotalNum += item.GoodsNum
		goodsId = append(goodsId, item.GoodId)
	}

	userInformation = models.GetUserInfoByOpenId(uid)
	goodsInfo = getCartListByIds(goodsId)

	cityId = userInformation.AddressDefault.CityId
	existAddress = !(len(userInformation.UserAddress) == 0)
	inRegion = true

	if len(carts) > 0 {
		for index, cart := range carts {
			if goodsInfo[cart.GoodId].GoodsId == 0 {
				carts = append(carts[:index], carts[index+1:]...)
				continue
			}
			good = goodsInfo[cart.GoodId]
			good.GoodsSkuId = cart.GoodsSkuId
			goodSku = GetGoodsSku(cart.GoodsSkuId, good)
			if goodSku.GoodsId == 0 {
				carts = append(carts[:index], carts[index+1:]...)
				continue
			}
			good.GoodsSku = goodSku
			if good.GoodsStatusArray["value"] != 10 {
				u.setErrorInfo(fmt.Sprintf("很抱歉，商品 [%s] 已下架", good.GoodsName))
			}

			if cart.GoodsNum > goodSku.StockNum {
				u.setErrorInfo(fmt.Sprintf("很抱歉，商品 [%s] 库存不足", good.GoodsName))
			}
			good.GoodsPrice = float64(goodSku.GoodsPrice)
			good.TotalNum = cart.GoodsNum
			good.TotalPrice = util.Multiplication(good.GoodsPrice * float64(cart.GoodsNum))
			good.GoodsTotalWeight = util.Multiplication(good.GoodsSku.GoodsWeight * float64(cart.GoodsNum))
			inRegion = checkAddress(cityId, good.Delivery.Rule)
			if inRegion {
				good.ExpressPrice = good.Delivery.GetTotalFee(cityId, cart.GoodsNum, good.GoodsTotalWeight)
			} else {
				if existAddress {
					u.setErrorInfo(fmt.Sprintf("很抱歉，您的收货地址不在商品 [%s] 的配送范围内", good.GoodsName))
				}
			}
			cartList = append(cartList, good)
		}
		orderTotalPrice, express = getTotalPriceAndExpress(cartList)
		expressPrice = getTotalExpressPrice(wxappId, express)
	}

	u.CartList = models.CartOrder{
		GoodList:        cartList,
		OrderTotalNum:   OrderTotalNum,
		OrderTotalPrice: util.Multiplication(orderTotalPrice),
		OrderPayPrice:   util.Multiplication(orderTotalPrice + expressPrice),
		Address:         userInformation.AddressDefault,
		ExistAddress:    existAddress,
		ExpressPrice:    expressPrice,
		IntraRegion:     inRegion,
		HasError:        u.hasError(),
		ErrorMsg:        u.getErrorInfo(),
	}
	return
}

//cartInfo := make(map[string]string)
//cartInfo["10002_1_3"] = "[{\"good_id\":10002,\"goods_num\":1,\"goods_sku_id\":\"1_3\",\"create_time\":1558496544}]"

func (u *UserCartList) Add(uid int) (err error) {
	var (
		key      string
		dataByte []byte
		exist    bool
		carts    []CartList
	)
	key = u.getKey(uid)
	if dataByte, exist, err = getDataFromRedis(key); err != nil {
		logging.LogError(err)
		return
	}

	if exist {
		if err = json.Unmarshal(dataByte, &carts); err != nil {
			logging.LogError(err)
			return
		}
		for index, item := range carts {
			if item.GoodId == u.AddGoodCart.GoodId {
				carts[index].GoodsNum++
				u.CartTotalNum = carts[index].GoodsNum
				carts[index].CreateTime = time.Now().Unix()
				break
			}
		}
	} else {
		u.CartTotalNum = u.AddGoodCart.GoodsNum
		carts = append(carts, u.AddGoodCart)
	}
	err = setDataWithKeyWithoutExpire(key, carts)
	return
}

//计算运费总结
func getTotalExpressPrice(wxappId string, express []float64) (expressPrice float64) {
	var (
		freightRule string
	)
	if len(express) == 0 {
		expressPrice = 0.00
		return
	}
	freightRule = models.GetSettingRuleId("trade", wxappId)

	switch freightRule {
	case "10":
		for _, item := range express {
			expressPrice += item
		}
	case "20":
		expressPrice = 0.00
		for _, item := range express {
			if item > expressPrice {
				expressPrice = item
			}
		}
	case "30":
		expressPrice = 0.00
		for _, item := range express {
			if item < expressPrice {
				expressPrice = item
			}
		}
	}
	return
}

//验证用户收货地址是否存在运费规则中
func checkAddress(cityId int, rules []models.DeliveryRule) (exists bool) {
	if cityId == 0 {
		return
	}

	regionId := strconv.Itoa(cityId)

	for _, item := range rules {
		if len(item.Region) == 0 {
			continue
		}
		if strings.Index(item.Region, regionId) > -1 {
			exists = true
			return
		}
	}
	return
}

//根据商品id集获取商品列表 (购物车列表用)
func getCartListByIds(goodsId []int) (goodsInfo map[int]models.Goods) {
	goods := models.GetGoodsInfoForCartList(goodsId)
	goodsInfo = make(map[int]models.Goods)
	for _, item := range goods {
		goodsInfo[item.GoodsId] = item
	}
	return
}

//获取商品总价和运费
func getTotalPriceAndExpress(goods []models.Goods) (price float64, express []float64) {
	for _, item := range goods {
		price += item.TotalPrice
		express = append(express, item.ExpressPrice)
	}
	return
}

// 商品多规格信息
func GetGoodsSku(goodSkuId string, g models.Goods) (goodSkuInfo models.GoodsSpec) {
	for _, item := range g.GoodsSpec {
		if item.SpecSkuId == goodSkuId {
			goodSkuInfo = item
		}
	}
	if goodSkuInfo.GoodsId == 0 {
		return
	}
	if g.SpecType == 20 {
		attrs := strings.Split(goodSkuInfo.SpecSkuId, "_")
		specRel := make(map[string]*models.SpecRel)
		goodsSpecRel, _ := GetGoodsSpecRel(g.GoodsId)
		g.SpecRel = goodsSpecRel
		for _, item := range goodsSpecRel {
			specRel[strconv.Itoa(item.SpecValueId)] = item
		}

		for _, attr := range attrs {
			specRelInfo := specRel[attr]
			goodSkuInfo.GoodsAttr += specRelInfo.Spec.SpecName + ":" + specRelInfo.SpecValue.SpecValue + ";"
		}
	}

	return
}
