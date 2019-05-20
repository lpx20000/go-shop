package models

import (
	"time"
)

type Categories []Category
type CategoriesWithChild []CategoryWithChild

type Category struct {
	CategoryId       uint       `json:"category_id"`
	Name             string     `json:"name"`
	ParentId         uint       `json:"parent_id"`
	ImageId          uint       `json:"image_id"`
	Sort             uint8      `json:"sort"`
	WxappId          uint       `json:"-"`
	CreateTime       int64      `json:"-"`
	CreateTimeString string     `json:"create_time"`
	Image            UploadFile `gorm:"foreignkey:ImageId;association_foreignkey:FileId" json:"image,omitempty" `
}

type CategoryWithChild struct {
	Category
	Child []Category `json:"child"`
}

func (c *Category) AfterFind() error {
	c.CreateTimeString = time.Unix(c.CreateTime, 0).Format("2006-01-02 15:04:05")
	return nil
}

func (c Categories) Len() int {
	return len(c)
}

func (c Categories) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Categories) Less(i, j int) bool {
	return c[i].Sort < c[j].Sort
}

func (c CategoriesWithChild) Len() int {
	return len(c)
}

func (c CategoriesWithChild) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CategoriesWithChild) Less(i, j int) bool {
	return c[i].Sort < c[j].Sort
}
