package models

type GoodsSpecRel struct {
	Id          uint `json:"id"`
	GoodsId     uint `json:"goods_id"`
	SpecId      uint `json:"spec_id"`
	SpecValueId uint `json:"spect_value_id"`
	WxappId     uint `json:"wxapp_id"`
}
