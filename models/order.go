package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	OrderId        int         `gorm:"PRIMARY_KEY" json:"order_id"`
	OrderNo        string      `json:"order_no"`
	TotalPrice     float64     `json:"total_price"`
	PayStatus      interface{} `gorm:"default:'10'" json:"pay_status"`
	PayPrice       float64     `json:"pay_price"`
	PayTime        uint        `json:"pay_time"`
	ExpressPrice   float64     `json:"express_price"`
	ExpressCompany string      `json:"express_company"`
	ExpressNo      string      `json:"express_no"`
	DeliveryStatus interface{} `gorm:"default:'10'" json:"delivery_status"`
	DeliveryTime   uint        `json:"delivery_time"`
	ReceiptStatus  interface{} `gorm:"default:'10'" json:"receipt_status"`
	ReceiptTime    uint        `json:"receipt_time"`
	OrderStatus    interface{} `gorm:"default:'10'"json:"order_status"`
	TransactionId  string      `json:"transaction_id"`
	UserId         int         `json:"user_id"`
	WxappId        string      `json:"-"`
	Model
	OrderGoods []OrderGoods `gorm:"foreignkey:order_id;association_foreignkey:order_id" json:"goods,omitempty" ` //hasMany
	//User               User                   `gorm:"foreignkey:UserId;association_foreignkey:UserId" json:"address,omitempty" ` //belongsTo
	OrderAddress OrderAddress `gorm:"foreignkey:OrderId;association_foreignkey:OrderId" json:"address,omitempty" `
}

type CartOrder struct {
	GoodList        []Goods     `json:"goods_list"`
	OrderTotalNum   int         `json:"order_total_num"`
	OrderTotalPrice float64     `json:"order_total_price"`
	OrderPayPrice   float64     `json:"order_pay_price"`
	Address         UserAddress `json:"address"`
	ExistAddress    bool        `json:"exist_address"`
	ExpressPrice    float64     `json:"express_price"`
	IntraRegion     bool        `json:"intra_region"`
	HasError        bool        `json:"has_error"`
	ErrorMsg        string      `json:"error_msg"`
}

func (o *Order) AfterFind() error {
	payStatus := map[int64]map[string]interface{}{
		10: {"text": "待付款", "value": 10},
		20: {"text": "已付款", "value": 20},
	}
	deliveryStatus := map[int64]map[string]interface{}{
		10: {"text": "待发货", "value": 10},
		20: {"text": "已发货'", "value": 20},
	}
	receiptStatus := map[int64]map[string]interface{}{
		10: {"text": "待发货", "value": 10},
		20: {"text": "已发货'", "value": 20},
	}
	orderStatus := map[int64]map[string]interface{}{
		10: {"text": "进行中", "value": 10},
		20: {"text": "取消'", "value": 20},
		30: {"text": "已完成'", "value": 30},
	}
	o.PayStatus = payStatus[o.PayStatus.(int64)]
	o.DeliveryStatus = deliveryStatus[o.DeliveryStatus.(int64)]
	o.ReceiptStatus = receiptStatus[o.ReceiptStatus.(int64)]
	o.OrderStatus = orderStatus[o.OrderStatus.(int64)]
	o.CreateTime, _ = strconv.ParseInt(time.Unix(o.CreateTime, 0).Format("2006-01-02 15:04:05"), 10, 64)
	return nil
}

func GetOrderList(userId int, filters map[string]interface{}) ([]*Order, error) {
	var (
		orders []*Order
		err    error
	)
	err = Db.Where(&Order{UserId: userId}).
		Where(filters).
		Preload("OrderGoods").
		Not("order_status", 20).
		Order("create_time DESC").
		Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return orders, nil
}

func GetOrderDetail(userId, orderId int) (*Order, error) {
	var (
		order Order
		err   error
	)
	err = Db.Where(&Order{UserId: userId, OrderId: orderId}).
		Preload("OrderGoods").
		Preload("OrderAddress").
		Not("order_status", 20).
		First(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &order, nil
}

func GetOrderCount(userId int, filter string) (count int) {
	var (
		filters map[string]interface{}
	)
	filters = make(map[string]interface{})
	switch filter {
	case "all":
	case "payment":
		filters["pay_status"] = 10
	case "received":
		filters["pay_status"] = 20
		filters["receipt_status"] = 20
		filters["receipt_status"] = 10
	}

	Db.Model(&Order{}).Where(&Order{UserId: userId}).
		Not("order_status", 20).
		Where(filters).Count(&count)
	return
}

func UpdateOrderByOrderId(orderId int, data map[string]interface{}) error {
	return Db.Model(&Order{}).Set("gorm:association_autoupdate", false).
		Where(&Order{OrderId: orderId}).Updates(data).Error
}
