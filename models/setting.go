package models

import (
	"encoding/json"
	"log"
)

type Setting struct {
	Key      string `json:"key"`
	Describe string `json:"describe"`
	Values   string `json:"values"`
	WxappId  int    `json:"wxapp_id"`
}

type Values struct {
	OrderValues OrderValues `json:"order"`
	FreightRule string      `json:"freight_rule"`
}

type OrderValues struct {
	CloseDays   string `json:"close_days"`
	ReceiveDays string `json:"receive_days"`
}

func GetSettingRuleId(key string, wxappId string) (ruleId string) {
	var (
		values  Values
		setting Setting
	)

	Db.Where(map[string]interface{}{"wxapp_id": wxappId, "key": key}).First(&setting)

	err := json.Unmarshal([]byte(setting.Values), &values)
	log.Println(err)
	ruleId = values.FreightRule
	return
}
