package models

import "time"

type GoodsSpecRel struct {
	GoodsSpecRefer
	GoodsId   uint      `json:"goods_id"`
	Spec      Spec      `gorm:"foreignkey:SpecId;association_foreignkey:SpecId" json:"spec" `       //belongsTo
	SpecValue SpecValue `gorm:"foreignkey:SpecId;association_foreignkey:SpecId" json:"spec_value" ` //belongsTo
}

type GoodsSpecRefer struct {
	Id              uint   `json:"id"`
	GoodsId         uint   `json:"goods_id"`
	SpecId          uint   `json:"spec_id"`
	SpecValueId     uint   `json:"spect_value_id"`
	WxappId         uint   `json:"-"`
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

func GetGoodsSpecRel(goodId uint) (specRelAll []SpecRel) {
	var goodsSpecRel []GoodsSpecRel
	db.Model(&GoodsSpecRel{}).
		Where(&GoodsSpecRel{GoodsId: goodId}).
		Preload("Spec").
		Preload("SpecValue").
		Find(&goodsSpecRel)

	for _, v := range goodsSpecRel {
		var spec SpecRel
		spec.SpecValue = v.SpecValue
		spec.Spec = v.Spec
		spec.Pivot.Id = v.Id
		spec.Pivot.GoodsId = v.GoodsId
		spec.Pivot.SpecId = v.SpecId
		spec.Pivot.SpecValueId = v.SpecValueId
		spec.Pivot.WxappId = v.WxappId
		spec.Pivot.CreateTimeStamp = v.CreateTimeStamp
		spec.Pivot.GoodsId = v.GoodsSpecRefer.GoodsId
		specRelAll = append(specRelAll, spec)
	}
	return
}
