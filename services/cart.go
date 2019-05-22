package services

import (
	"encoding/json"
	"fmt"
	"shop/models"
)

type cartList struct {
	GoodID     uint   `json:"good_id"`
	GoodsNum   uint   `json:"goods_num"`
	GoodsSkuID string `json:"goods_sku_id"`
	CreateTime int    `json:"create_time"`
}

type ErrorInfo struct {
	err string
}

func (e *ErrorInfo) setErrorInfo(err string) {
	e.err = err
}

func (e *ErrorInfo) getErrorInfo() (err string) {
	return e.err
}

func GetCartInfo(token interface{}) {
	var (
		carts           []cartList
		goodsId         []uint
		goodsInfo       map[uint]models.Goods
		userInformation models.User
		cityId          uint
		existAddress    bool
		intraRegion     bool
		cartList        []models.Goods
		good            models.Goods
		goodSku         models.GoodsSpec
		errObject       ErrorInfo
	)
	cartInfo := make(map[string]string)
	cartInfo["10002_1_3"] = "[{\"good_id\":10002,\"goods_num\":1,\"goods_sku_id\":\"1_3\",\"create_time\":1558496544}]"

	_ = json.Unmarshal([]byte(cartInfo["10002_1_3"]), &carts)

	for _, item := range carts {
		goodsId = append(goodsId, item.GoodID)
	}

	userInformation = models.GetUserInfoByOpenId(token)
	goodsInfo = getCartListByIds(goodsId)

	cityId = userInformation.AddressDefault.CityId
	existAddress = len(userInformation.UserAddress) == 0
	intraRegion = true

	for index, cart := range carts {
		if goodsInfo[cart.GoodID].GoodsId == 0 {
			carts = append(carts[:index], carts[index+1:]...)
			continue
		}
		good = goodsInfo[cart.GoodID]
		good.GoodsSkuId = cart.GoodsSkuID
		goodSku = getGoodsSku(cart.GoodsSkuID, goodsInfo[cart.GoodID].GoodsSpec)
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
		good.GoodsPrice = goodSku.GoodsPrice
		good.TotalNum = cart.GoodsNum
		totalPrice := good.GoodsPrice * float32(cart.GoodsNum)
		good.TotalPrice = fmt.Sprintf("%0.2f", totalPrice)
	}

}

func getCartListByIds(goodsId []uint) (goodsInfo map[uint]models.Goods) {
	goods := models.GetGoodsInfoForCartList(goodsId)
	goodsInfo = make(map[uint]models.Goods)
	for _, item := range goods {
		goodsInfo[item.GoodsId] = item
	}
	return
}

func getGoodsSku(goodSkuId string, goodsSpec []models.GoodsSpec) models.GoodsSpec {

}
