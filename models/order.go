package models

type Order struct {
	OrderId        uint    `json:"order_id"`
	OrderNo        string  `json:"order_no"`
	TotalPrice     float32 `json:"total_price"`
	PayNotice      float32 `json:"pay_notice"`
	PayStatus      uint8   `json:"pay_status"`
	PayTime        uint    `json:"pay_time"`
	ExpressPrice   float32 `json:"express_price"`
	ExpressCompany string  `json:"express_company"`
	ExpressNo      string  `json:"express_no"`
	DeliveryStatus uint8   `json:"delivery_status"`
	DeliveryTime   uint    `json:"delivery_time"`
	ReceiptTime    uint    `json:"receipt_time"`
	OrderStatus    uint8   `json:"order_status"`
	TransactionId  string  `json:"transaction_id"`
	UserId         uint    `json:"user_id"`
	WxappId        uint    `json:"wxapp_id"`
}

type CartOrder struct {
	GoodList        []Goods     `json:"good_list"`
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
