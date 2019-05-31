package models

type Order struct {
	OrderId            uint                   `json:"order_id"`
	OrderNo            string                 `json:"order_no"`
	TotalPrice         float32                `json:"total_price"`
	PayNotice          float32                `json:"pay_notice,omitempty"`
	PayStatus          uint8                  `json:"-"`
	PayStatusInfo      map[string]interface{} `json:"pay_status,omitempty"`
	PayTime            uint                   `json:"pay_time"`
	ExpressPrice       float32                `json:"express_price"`
	ExpressCompany     string                 `json:"express_company"`
	ExpressNo          string                 `json:"express_no"`
	DeliveryStatus     uint8                  `json:"-"`
	DeliveryStatusInfo map[string]interface{} `json:"delivery_status,omitempty"`
	DeliveryTime       uint                   `json:"delivery_time"`
	ReceiptStatus      uint8                  `json:"-"`
	ReceiptStatusInfo  map[string]interface{} `json:"receipt_status,omitempty"`
	ReceiptTime        uint                   `json:"receipt_time"`
	OrderStatus        uint8                  `json:"-"`
	OrderStatusInfo    map[string]interface{} `json:"order_status,omitempty"`
	TransactionId      string                 `json:"transaction_id"`
	UserId             int                    `json:"user_id"`
	WxappId            uint                   `json:"-"`
	OrderGoods         []OrderGoods           `gorm:"foreignkey:OrderId;association_foreignkey:OrderId" json:"goods,omitempty" ` //hasMany
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
	payStatus := map[uint8]map[string]interface{}{
		10: {"text": "待付款", "value": 10},
		20: {"text": "已付款", "value": 20},
	}
	deliveryStatus := map[uint8]map[string]interface{}{
		10: {"text": "待发货", "value": 10},
		20: {"text": "已发货'", "value": 20},
	}
	receiptStatus := map[uint8]map[string]interface{}{
		10: {"text": "待发货", "value": 10},
		20: {"text": "已发货'", "value": 20},
	}
	orderStatus := map[uint8]map[string]interface{}{
		10: {"text": "进行中", "value": 10},
		20: {"text": "取消'", "value": 20},
		30: {"text": "已完成'", "value": 30},
	}
	o.PayStatusInfo = payStatus[o.PayStatus]
	o.DeliveryStatusInfo = deliveryStatus[o.DeliveryStatus]
	o.ReceiptStatusInfo = receiptStatus[o.ReceiptStatus]
	o.OrderStatusInfo = orderStatus[o.OrderStatus]
	return nil
}

func GetOrderList(userId int, filter string) (orders []Order) {
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

	Db.Where(&Order{UserId: userId}).
		Where(filters).
		Preload("OrderGoods").
		Not("order_status", 20).
		Order("create_time DESC").
		Find(&orders)
	return
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
		Preload("GoodsImage").
		Not("order_status", 20).
		Order("create_time DESC").
		Where(filters).Count(&count)
	return
}
