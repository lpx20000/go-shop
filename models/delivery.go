package models

type Delivery struct {
	DeliveryId uint   `json:"delivery_id"`
	Name       string `json:"name"`
	Method     uint   `json:"method"`
	Sort       uint8  `json:"sort"`
	WxappId    uint   `json:"wxapp_id"`
}
