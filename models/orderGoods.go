package models

import "github.com/jinzhu/gorm"

type OrderGoods struct {
	CreateTime      int64      `json:"-"`
	OrderGoodsId    int        `json:"order_goods_id"`
	OrderId         int        `json:"order_id"`
	GoodsId         int        `json:"goods_id"`
	GoodsName       string     `json:"goods_name"`
	GoodsSpecId     int        `json:"goods_spec_id"`
	ImageId         uint       `json:"image_id"`
	DeductStockType int        `json:"deduct_stock_type"`
	SpecType        int        `json:"spec_type"`
	SpecSkuId       string     `json:"spec_sku_id"`
	GoodsAttr       string     `json:"goods_attr"`
	Content         string     `json:"content,omitempty"`
	GoodsNo         string     `json:"goods_no"`
	GoodsPrice      float64    `json:"goods_price"`
	LinePrice       float32    `json:"line_price"`
	GoodsWeight     float64    `json:"goods_weight"`
	TotalNum        int        `json:"total_num"`
	TotalPrice      float64    `json:"total_price"`
	UserId          int        `json:"user_id"`
	WxappId         string     `json:"-"`
	Image           UploadFile `gorm:"foreignkey:ImageId;association_foreignkey:FileId" json:"image,omitempty" `         //belongsTo
	Goods           Goods      `gorm:"foreignkey:GoodsId;association_foreignkey:GoodsId" json:"goods,omitempty" `        //belongsTo
	GoodsSpec       GoodsSpec  `gorm:"foreignkey:GoodsSpecId;association_foreignkey:GoodsSpecId" json:"spec,omitempty" ` //belongsTo
}

func GetOrderGoods(orderId int) (orderGoods []OrderGoods, err error) {
	err = Db.Model(&OrderGoods{}).
		Preload("Image").
		Preload("Goods").
		Preload("GoodsSpec").
		Where(&OrderGoods{OrderId: orderId}).Find(&orderGoods).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}
