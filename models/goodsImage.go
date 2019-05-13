package models

type GoodsImage struct {
	Id      uint `json:"id"`
	GoodsId uint `json:"goods_id"`
	ImageId uint `json:"image_id"`
	WxappId uint `json:"wxapp_id"`
}
