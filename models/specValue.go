package models

type SpecValue struct {
	SpecValueId uint   `json:"spec_value_id"`
	SpecValue   string `json:"spec_value"`
	SpecId      uint   `json:"spec_id"`
	WxappId     uint   `json:"wxapp_id"`
}
