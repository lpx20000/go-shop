package services

import (
	"errors"
	"fmt"
	"math/rand"
	"shop/models"
	"shop/pkg/logging"
	"shop/pkg/util"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/chanxuehong/wechat/mch/core"
	"github.com/chanxuehong/wechat/mch/pay"
)

type Order struct {
	OrderList   OrderList   `json:"-"`
	UserId      int         `json:"-"`
	WxappId     string      `json:"wxapp_id"`
	ClientIp    string      `json:"-"`
	Token       string      `json:"-"`
	OrderDetail OrderDetail `json:"-"`
	Err         ErrorInfo   `json:"-"`
	OrderBuy    OrderBuy    `json:"-"`
	Mchid       string      `json:"-"`
	AppId       string      `json:"-"`
	AppSecret   string      `json:"-"`
	ApiKey      string      `json:"-"`
}

type OrderBuy struct {
	GoodsId    int              `form:"goods_id" binding:"required" json:"-"`
	GoodsNum   int              `form:"goods_num" binding:"required" json:"-"`
	GoodsSkuId string           `form:"goods_sku_id"  json:"-"`
	CartOrder  models.CartOrder `json:"-"`
	WxPay      WxPay            `json:"-"`
}

type WxPay struct {
	Payment   string `json:"-"`
	OrderId   int    `json:"-"`
	OrderNo   string `json:"-"`
	PrepayId  string `json:"prepay_id"`
	NonceStr  string `json:"nonce_str"`
	TimeStamp string `json:"timeStamp"`
	PaySign   string `json:"paySign"`
}

type OrderList struct {
	DataType string          `form:"dataType" binding:"required" json:"-"`
	List     []*models.Order `json:"list"`
}

type OrderDetail struct {
	OrderId int           `form:"order_id" binding:"required" json:"-"`
	Detail  *models.Order `json:"order"`
}

func (o *Order) GetOrderList(filter string) (err error) {
	var (
		filters map[string]interface{}
	)
	filters = make(map[string]interface{})
	switch filter {
	case "all":
	case "payment":
		filters["pay_status"] = 10
	case "delivery":
		filters["pay_status"] = 20
		filters["delivery_status"] = 10
	case "received":
		filters["pay_status"] = 20
		filters["receipt_status"] = 20
		filters["receipt_status"] = 10
	}
	if o.OrderList.List, err = models.GetOrderList(o.UserId, filters); err != nil {
		logging.LogTrace(err)
		return
	}
	return
}

func (o *Order) GetOrderDetail() (err error) {
	if err = o.GetSpecifyOrder(); err != nil {
		logging.LogTrace(err)
		return
	}
	if o.OrderDetail.Detail.OrderGoods, err = models.GetOrderGoods(o.OrderDetail.Detail.OrderId); err != nil {
		logging.LogTrace(err)
		return
	}
	return
}

//获取指定订单
func (o *Order) GetSpecifyOrder() (err error) {
	o.OrderDetail.Detail, err = models.GetOrderDetail(o.UserId, o.OrderDetail.OrderId)
	return
}

//订单取消
func (o *Order) CancelOrder() (err error) {
	if err = o.GetSpecifyOrder(); err != nil {
		logging.LogTrace(err)
		return
	}
	if o.OrderDetail.Detail.PayStatus == 20 {
		err = errors.New("已付款订单不可取消")
		return
	}
	tx := models.Db.Begin()
	if err = o.backGoodsStock(tx); err != nil {
		logging.LogTrace(err)
		tx.Rollback()
		return
	}
	err = tx.Model(&models.Order{}).
		Where(&models.Order{OrderId: o.OrderDetail.OrderId}).
		UpdateColumn("order_status", 20).
		Error
	if err != nil {
		tx.Rollback()
		logging.LogTrace(err)
	}
	tx.Commit()
	return
}

//确认收货
func (o *Order) ReceiptOrder() (err error) {
	if err = o.GetSpecifyOrder(); err != nil {
		logging.LogTrace(err)
		return
	}
	if o.OrderDetail.Detail.DeliveryStatus == 10 || o.OrderDetail.Detail.ReceiptStatus == 20 {
		err = errors.New("该订单不合法")
		return
	}
	data := make(map[string]interface{}, 3)
	data["receipt_status"] = 20
	data["receipt_time"] = time.Now().Unix()
	data["order_status"] = 30

	if err = models.UpdateOrderByOrderId(o.OrderDetail.OrderId, data); err != nil {
		logging.LogTrace(err)
	}
	return
}

//获取立即购买商品详细信息
func (o *Order) DoBuyToCart() (err error) {
	if err = o.getBuyNow(); err != nil || o.Err.HasError() {
		if len(o.Err.GetErrorInfo()) > 0 {
			logging.LogTrace(err)
			return
		}
		logging.LogTrace(err)
		return
	}
	return
}

//提交订单付款
func (o *Order) OrderBuyNow() (err error) {
	var order *models.Order
	if order, err = o.addOrder(); err != nil {
		logging.LogTrace(err)
		return
	}
	if err = o.getWxPayInfo(); err != nil {
		logging.LogTrace(err)
		return
	}
	o.OrderBuy.WxPay.OrderId = order.OrderId
	o.OrderBuy.WxPay.Payment = fmt.Sprintf("%.f", o.OrderBuy.CartOrder.OrderPayPrice*100)

	return
}

func (o *Order) WxPay() (err error) {
	//处理微信支付
	if err = o.initWxPay(); err != nil {
		logging.LogTrace(err)
	}
	return
}

//获取支付配置信息
func (o *Order) getWxPayInfo() (err error) {
	var appInfo models.Wxapp
	if appInfo, err = models.GetAppPayInfo(o.WxappId); err != nil {
		return
	}
	o.AppId = appInfo.AppId
	o.Mchid = appInfo.Mchid
	o.ApiKey = appInfo.Apikey
	o.AppSecret = appInfo.AppSecret
	return
}

//添加订单到数据库
func (o *Order) addOrder() (order *models.Order, err error) {
	if o.OrderBuy.CartOrder.Address.AddressId == 0 {
		err = errors.New("【添加订单】请先选择收货地址")
		return
	}
	o.OrderBuy.WxPay.OrderNo = o.getOrderNo()
	order = &models.Order{
		UserId:       o.UserId,
		WxappId:      o.WxappId,
		OrderNo:      o.OrderBuy.WxPay.OrderNo,
		TotalPrice:   o.OrderBuy.CartOrder.OrderTotalPrice,
		PayPrice:     o.OrderBuy.CartOrder.OrderPayPrice,
		ExpressPrice: o.OrderBuy.CartOrder.ExpressPrice,
		OrderGoods:   o.addOrderGoods(),
		OrderAddress: o.addOrderAddress(),
	}
	err = models.Db.Create(order).Error
	return
}

//添加订单商品
func (o *Order) addOrderGoods() (orderGood []models.OrderGoods) {
	for _, goodList := range o.OrderBuy.CartOrder.GoodList {
		goodOrder := models.OrderGoods{
			UserId:          o.UserId,
			OrderId:         o.OrderBuy.WxPay.OrderId,
			WxappId:         o.WxappId,
			GoodsId:         goodList.GoodsId,
			GoodsName:       goodList.GoodsName,
			ImageId:         goodList.GoodsImage[0].ImageId,
			DeductStockType: goodList.DeductStockType,
			SpecType:        goodList.SpecType,
			SpecSkuId:       goodList.GoodsSku.SpecSkuId,
			GoodsSpecId:     goodList.GoodsSku.GoodsSpecId,
			GoodsAttr:       goodList.GoodsSku.GoodsAttr,
			Content:         goodList.Content,
			GoodsNo:         goodList.GoodsSku.GoodsNo,
			GoodsPrice:      goodList.GoodsSku.GoodsPrice,
			LinePrice:       goodList.GoodsSku.LinePrice,
			GoodsWeight:     goodList.GoodsSku.GoodsWeight,
			TotalNum:        goodList.TotalNum,
			TotalPrice:      goodList.TotalPrice,
		}
		if goodList.DeductStockType == 10 {
			goodOrder.GoodsSpec = models.GoodsSpec{
				GoodsSpecId: goodList.GoodsSku.GoodsSpecId,
				StockNum:    -goodList.TotalNum,
			}
		}
		orderGood = append(orderGood, goodOrder)
	}
	return
}

//添加订单地址
func (o *Order) addOrderAddress() (address models.OrderAddress) {
	address = models.OrderAddress{
		UserId:     o.UserId,
		WxappId:    o.WxappId,
		Name:       o.OrderBuy.CartOrder.Address.Name,
		Phone:      o.OrderBuy.CartOrder.Address.Phone,
		ProvinceId: o.OrderBuy.CartOrder.Address.ProvinceId,
		CityId:     o.OrderBuy.CartOrder.Address.CityId,
		RegionId:   o.OrderBuy.CartOrder.Address.RegionId,
		Detail:     o.OrderBuy.CartOrder.Address.Detail,
	}
	return
}

//获取立即购买商品信息
func (o *Order) getBuyNow() (err error) {
	var (
		good            models.Goods
		user            models.User
		existAddress    bool
		inRegion        bool
		orderTotalPrice float64
		expressPrice    float64
		express         []float64
		goods           []models.Goods
	)

	if good, err = models.GetGoodDetail(o.OrderBuy.GoodsId); err != nil {
		logging.LogTrace(err)
		return
	}

	if good.GoodsStatus != 10 {
		o.Err.SetErrorInfo(fmt.Sprintf("很抱歉，商品 [%s] 已下架", good.GoodsName))
	}
	good.GoodsSku = GetGoodsSku(o.OrderBuy.GoodsSkuId, good)

	if o.OrderBuy.GoodsNum > good.GoodsSku.StockNum {
		o.Err.SetErrorInfo(fmt.Sprintf("很抱歉，商品 [%s] 库存不足", good.GoodsName))
	}
	good.GoodsPrice = good.GoodsSku.GoodsPrice
	good.TotalNum = o.OrderBuy.GoodsNum
	good.TotalPrice = util.Multiplication(good.GoodsPrice * float64(o.OrderBuy.GoodsNum))
	good.GoodsTotalWeight = util.Multiplication(good.GoodsSku.GoodsWeight * float64(o.OrderBuy.GoodsNum))
	user = models.GetUserInfoByOpenId(o.UserId)
	existAddress = !(len(user.UserAddress) == 0)
	inRegion = checkAddress(user.AddressDefault.CityId, good.Delivery.Rule)
	if inRegion {
		good.ExpressPrice = good.Delivery.GetTotalFee(user.AddressDefault.CityId,
			o.OrderBuy.GoodsNum, good.GoodsTotalWeight)
	} else {
		if existAddress {
			o.Err.SetErrorInfo(fmt.Sprintf("很抱歉，您的收货地址不在商品 [%s] 的配送范围内", good.GoodsName))
		}
	}

	goods = append(goods, good)
	orderTotalPrice, express = getTotalPriceAndExpress(goods)
	expressPrice = getTotalExpressPrice(o.WxappId, express)
	o.OrderBuy.CartOrder = models.CartOrder{
		GoodList:        goods,
		OrderTotalNum:   o.OrderBuy.GoodsNum,
		OrderTotalPrice: util.Multiplication(orderTotalPrice),
		OrderPayPrice:   util.Multiplication(orderTotalPrice + expressPrice),
		Address:         user.AddressDefault,
		ExistAddress:    existAddress,
		ExpressPrice:    expressPrice,
		IntraRegion:     inRegion,
		HasError:        o.Err.HasError(),
		ErrorMsg:        o.Err.GetErrorInfo(),
	}
	return
}

//生成订单好
func (o *Order) getOrderNo() (orderNo string) {
	orderNo = fmt.Sprintf("%s%08v", time.Unix(time.Now().Unix(), 0).Format("20060102"), rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
	return
}

//取消订单返回库存
func (o *Order) backGoodsStock(tx *gorm.DB) (err error) {
	for _, goodList := range o.OrderDetail.Detail.OrderGoods {
		if goodList.DeductStockType == 10 {
			err = tx.Model(&models.GoodsSpec{}).
				Where(map[string]interface{}{"goods_spec_id": goodList.GoodsSpecId}).
				UpdateColumn("stock_num", gorm.Expr("stock_num + ?", goodList.TotalNum)).Error
			if err != nil {
				return
			}
		}
	}
	return
}

//初始化微信支付
func (o *Order) initWxPay() (err error) {
	var (
		params    map[string]string
		resp      map[string]string
		paramsPre map[string]string
		timeStamp string
		monceStr  string
	)
	timeStamp = fmt.Sprintf("%d", time.Now().Unix())
	monceStr = util.Md5V(fmt.Sprintf("%s%s", timeStamp, o.OrderBuy.WxPay.OrderNo))
	params = make(map[string]string, 11)
	params["appid"] = o.AppId
	params["attach"] = fmt.Sprintf("%s%s", o.WxappId, "支付")
	params["body"] = o.OrderBuy.WxPay.OrderNo
	params["mch_id"] = o.Mchid
	params["nonce_str"] = monceStr
	params["notify_url"] = "notice"
	params["openid"] = o.Token
	params["out_trade_no"] = o.OrderBuy.WxPay.OrderNo
	params["spbill_create_ip"] = o.ClientIp
	params["total_fee"] = o.OrderBuy.WxPay.Payment
	params["trade_type"] = "JSAPI"

	resp, err = pay.UnifiedOrder(core.NewClient(o.AppId, o.Mchid, o.ApiKey, nil), params)
	if err != nil {
		return
	}

	if resp["return_code"] == "FAIL" {
		err = errors.New(resp["return_msg"])
		return
	}

	paramsPre = make(map[string]string, 5)
	paramsPre["appid"] = o.AppId
	paramsPre["nonceStr"] = monceStr
	paramsPre["package"] = fmt.Sprintf("%s%s", "prepay_id=", resp["prepay_id"])
	paramsPre["signType"] = "MD5"
	paramsPre["timeStamp"] = timeStamp

	o.OrderBuy.WxPay.PrepayId = resp["prepay_id"]
	o.OrderBuy.WxPay.NonceStr = monceStr
	o.OrderBuy.WxPay.TimeStamp = timeStamp
	o.OrderBuy.WxPay.PaySign = core.Sign(paramsPre, o.ApiKey, nil)
	return
}
