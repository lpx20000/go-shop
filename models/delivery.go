package models

import (
	"fmt"
	"shop/pkg/util"
	"strconv"
	"strings"
)

type Delivery struct {
	DeliveryId  uint                   `json:"delivery_id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Method      uint                   `json:"-"`
	MethodSlice map[string]interface{} `json:"method,omitempty"`
	Sort        uint8                  `json:"sort,omitempty"`
	WxappId     uint                   `json:"-"`
	Rule        []DeliveryRule         `gorm:"foreignkey:DeliveryId;association_foreignkey:DeliveryId" json:"rule,omitempty" `
}

func (d *Delivery) AfterFind() error {
	Db.Model(&d).Related(&d.Rule, "Rule")
	goodsStatus := map[uint]map[string]interface{}{
		10: {"text": "按件数", "value": 10},
		20: {"text": "'按重量", "value": 20},
	}
	d.MethodSlice = goodsStatus[d.Method]
	return nil
}

func (d *Delivery) GetTotalFee(cityId, goodsNum int, goodWeight float64) (fee float64) {
	var (
		ruleInfo      DeliveryRule
		total         float64
		additional    float64
		additionalFee float64
	)

	total = goodWeight
	regionId := strconv.Itoa(cityId)
	for _, item := range d.Rule {
		if strings.Index(item.Region, regionId) > -1 {
			ruleInfo = item
			break
		}
	}
	if d.MethodSlice["value"] == 10 {
		total = float64(goodsNum)
	}
	additional = total - ruleInfo.First

	if additional <= 0 {
		fee = ruleInfo.FirstFee
		return
	}

	if additional <= ruleInfo.Additional {
		fee, _ = util.FormatNumber(fmt.Sprintf("%.2f", ruleInfo.FirstFee+ruleInfo.AdditionalFee))
		return
	}

	additionalFee = 0.00
	if ruleInfo.Additional >= 1 {
		additionalFee = util.Multiplication(ruleInfo.AdditionalFee/ruleInfo.Additional) * additional
	}
	fee = util.Multiplication(ruleInfo.FirstFee + additionalFee)
	return
}
