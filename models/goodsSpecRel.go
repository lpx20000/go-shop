package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type GoodsSpecRel struct {
	GoodsSpecRefer
	GoodsId   int       `json:"goods_id"`
	Spec      Spec      `gorm:"foreignkey:SpecId;association_foreignkey:SpecId" json:"spec" `                 //belongsTo
	SpecValue SpecValue `gorm:"foreignkey:SpecValueId;association_foreignkey:SpecValueId" json:"spec_value" ` //belongsTo
}

type GoodsSpecRefer struct {
	Id              uint   `json:"id"`
	GoodsId         int    `json:"goods_id"`
	SpecId          int    `json:"spec_id"`
	SpecValueId     int    `json:"spect_value_id"`
	WxappId         int    `json:"-"`
	CreateTime      int64  `json:"-"`
	CreateTimeStamp string `json:"create_time"`
}

type SpecRel struct {
	SpecValue
	Spec  Spec           `json:"spec"`
	Pivot GoodsSpecRefer `json:"pivot"`
}

func (g *GoodsSpecRel) AfterFind() error {
	g.SpecValue.CreateTimeStamp = time.Unix(g.SpecValue.CreateTime, 0).Format("2006-01-02 15:04:05")
	g.Spec.CreateTimeStamp = time.Unix(g.Spec.CreateTime, 0).Format("2006-01-02 15:04:05")
	g.CreateTimeStamp = time.Unix(g.CreateTime, 0).Format("2006-01-02 15:04:05")
	return nil
}

func GetGoodSpecRel(goodId int) ([]*GoodsSpecRel, error) {
	var (
		goodsSpecRel []*GoodsSpecRel
		err          error
	)
	err = Db.Model(&GoodsSpecRel{}).
		Where(&GoodsSpecRel{GoodsId: goodId}).
		Preload("Spec").
		Preload("SpecValue").
		Find(&goodsSpecRel).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return goodsSpecRel, nil
}
