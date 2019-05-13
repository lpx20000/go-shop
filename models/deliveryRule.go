package models

type DeliveryRule struct {
	RuleId        uint    `json:"rule_id"`
	DeliveryId    uint    `json:"delivery_id"`
	Region        string  `json:"region"`
	First         float64 `json:"first"`
	FirstFee      float32 `json:"first_fee"`
	Additional    float64 `json:"additional"`
	AdditionalFee float32 `json:"additional_fee"`
	WxappId       uint    `json:"wxapp_id"`
}
