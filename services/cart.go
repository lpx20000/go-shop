package services

import (
	"encoding/json"
	"fmt"
	"shop/models"
	"shop/pkg/util"
	"strconv"
	"strings"
)

type cartList struct {
	GoodID     uint   `json:"good_id"`
	GoodsNum   int    `json:"goods_num"`
	GoodsSkuID string `json:"goods_sku_id"`
	CreateTime int    `json:"create_time"`
}

type ErrorInfo struct {
	err    string
	exists bool
}

func (e *ErrorInfo) setErrorInfo(err string) {
	e.exists = true
	e.err = err
}

func (e *ErrorInfo) getErrorInfo() (err string) {
	err = e.err
	e.exists = false
	e.err = ""
	return
}

func (e *ErrorInfo) hasError() bool {
	return e.exists
}

func GetCartInfo(token string, wxappId string) (cart models.CartOrder) {
	var (
		carts           []cartList
		goodsId         []uint
		goodsInfo       map[uint]models.Goods
		userInformation models.User
		cityId          int
		existAddress    bool
		cartList        []models.Goods
		good            models.Goods
		goodSku         models.GoodsSpec
		errObject       ErrorInfo
		inRegion        bool
		orderTotalPrice float64
		expressPrice    float64
		express         []float64
		OrderTotalNum   int
	)
	cartInfo := make(map[string]string)
	cartInfo["10002_1_3"] = "[{\"good_id\":10002,\"goods_num\":1,\"goods_sku_id\":\"1_3\",\"create_time\":1558496544}]"

	_ = json.Unmarshal([]byte(cartInfo["10002_1_3"]), &carts)

	for _, item := range carts {
		OrderTotalNum += item.GoodsNum
		goodsId = append(goodsId, item.GoodID)
	}

	userInformation = models.GetUserInfoByOpenId(token)
	goodsInfo = getCartListByIds(goodsId)

	cityId = userInformation.AddressDefault.CityId
	existAddress = !(len(userInformation.UserAddress) == 0)
	inRegion = true

	for index, cart := range carts {
		if goodsInfo[cart.GoodID].GoodsId == 0 {
			carts = append(carts[:index], carts[index+1:]...)
			continue
		}
		good = goodsInfo[cart.GoodID]
		good.GoodsSkuId = cart.GoodsSkuID
		goodSku = GetGoodsSku(cart.GoodsSkuID, good)
		if goodSku.GoodsId == 0 {
			carts = append(carts[:index], carts[index+1:]...)
			continue
		}
		good.GoodsSku = goodSku
		if good.GoodsStatusArray["value"] != 10 {
			errObject.setErrorInfo(fmt.Sprintf("很抱歉，商品 [%s] 已下架", good.GoodsName))
		}

		if cart.GoodsNum > goodSku.StockNum {
			errObject.setErrorInfo(fmt.Sprintf("很抱歉，商品 [%s] 库存不足", good.GoodsName))
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
				errObject.setErrorInfo(fmt.Sprintf("很抱歉，您的收货地址不在商品 [%s] 的配送范围内", good.GoodsName))
			}
		}
		cartList = append(cartList, good)
	}
	orderTotalPrice, express = getTotalPriceAndExpress(cartList)
	expressPrice = getTotalExpressPrice(wxappId, express)

	cart = models.CartOrder{
		GoodList:        cartList,
		OrderTotalNum:   OrderTotalNum,
		OrderTotalPrice: util.Multiplication(orderTotalPrice),
		OrderPayPrice:   util.Multiplication(orderTotalPrice + expressPrice),
		Address:         userInformation.AddressDefault,
		ExistAddress:    existAddress,
		ExpressPrice:    expressPrice,
		IntraRegion:     inRegion,
		HasError:        errObject.hasError(),
		ErrorMsg:        errObject.getErrorInfo(),
	}

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
func getCartListByIds(goodsId []uint) (goodsInfo map[uint]models.Goods) {
	goods := models.GetGoodsInfoForCartList(goodsId)
	goodsInfo = make(map[uint]models.Goods)
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
		specRel := make(map[string]models.SpecRel)
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

func GetRegionInfo() (tree map[int]models.Tree) {
	var (
		commonList []models.CommonList
	)

	tree = make(map[int]models.Tree)

	commonList = models.GetRegion()

	for _, province := range commonList {

		if province.Level == 1 {
			tree[province.Id] = models.Tree{
				CommonList: province,
				City:       make(map[int]models.City, 0),
			}

			for _, city := range commonList {
				if city.Level == 2 && city.Pid == province.Id {
					tree[province.Id].City[city.Id] = models.City{
						CommonList: city,
						RegionInfo: make(map[int]models.CommonList, 0),
					}

					for _, region := range commonList {
						if region.Level == 3 && region.Pid == city.Id {
							tree[province.Id].City[city.Id].RegionInfo[region.Id] = region
						}
					}
				}
			}
		}
	}
	return
}
