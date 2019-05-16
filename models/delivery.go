package models

type Delivery struct {
	DeliveryId  uint                   `json:"delivery_id"`
	Name        string                 `json:"name"`
	Method      uint                   `json:"-"`
	MethodSlice map[string]interface{} `json:"method"`
	Sort        uint8                  `json:"sort"`
	WxappId     uint                   `json:"-"`
	Rule        []DeliveryRule         `gorm:"foreignkey:DeliveryId;association_foreignkey:DeliveryId" json:"rule" `
}

func (d *Delivery) AfterFind() error {
	db.Model(&d).Related(&d.Rule, "Rule")
	goodsStatus := map[uint]map[string]interface{}{
		10: {"text": "按件数", "value": 10},
		20: {"text": "'按重量", "value": 20},
	}
	d.MethodSlice = goodsStatus[d.Method]
	return nil
}
