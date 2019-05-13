package models

type Spec struct {
	SpecId   uint   `json:"spec_id"`
	SpecName string `json:"spec_name"`
	WxappId  uint   `json:"wxapp_id"`
}
