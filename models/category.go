package models

import "time"

type Category struct {
	CategoryId       uint   `json:"category_id"`
	Name             string `json:"name"`
	ParentId         uint   `json:"parent_id"`
	ImageId          uint   `json:"image_id"`
	Sort             uint8  `json:"sort"`
	WxappId          uint   `json:"-"`
	CreateTime       int64  `json:"-"`
	CreateTimeString string `json:"create_time"`
}

func (c *Category) AfterFind() error {
	c.CreateTimeString = time.Unix(c.CreateTime, 0).Format("2006-01-02 15:04:05")
	return nil
}
