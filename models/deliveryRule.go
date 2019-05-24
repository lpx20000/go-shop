package models

import "strings"

type DeliveryRule struct {
	RuleId        uint     `json:"rule_id"`
	DeliveryId    uint     `json:"delivery_id"`
	Region        string   `json:"region"`
	RegionData    []string `json:"region_data"`
	First         float64  `json:"first"`
	FirstFee      float64  `json:"first_fee"`
	Additional    float64  `json:"additional"`
	AdditionalFee float64  `json:"additional_fee"`
	WxappId       uint     `json:"-"`
}

func (d *DeliveryRule) AfterFind() error {
	d.RegionData = strings.Split(d.Region, ",")
	return nil
}
