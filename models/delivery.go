package models

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
