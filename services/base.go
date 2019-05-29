package services

import (
	"encoding/json"
	"shop/pkg/gredis"
	"shop/pkg/logging"
)

type Base struct {
}

func (b *Base) GetDataFromRedis(key string) (exist bool, err error) {
	if gredis.Exists(key) {
		var (
			dataByte []byte
		)
		if dataByte, err = gredis.Get(key); err != nil {
			return
		}
		if err = json.Unmarshal(dataByte, b); err != nil {
			logging.LogInfo(err.Error())
			return
		}
		exist = true
	}
	return
}

func (b *Base) SetDataWithKey(key string, data interface{}) error {
	return gredis.Set(key, data, 3600)
}
