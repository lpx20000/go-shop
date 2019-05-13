package models

type Category struct {
	CategoryId uint   `json:"category_id"`
	Name       string `json:"name"`
	ParentId   uint   `json:"parent_id"`
	ImageId    uint   `json:"image_id"`
	Sort       uint8  `json:"sort"`
	WxappId    uint   `json:"wxapp_id"`
}
