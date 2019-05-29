package services

import (
	"shop/pkg/gredis"
)

type Base struct {
}

func (h *Base) getDataFromRedis(key string) (dataByte []byte, exist bool, err error) {
	if gredis.Exists(key) {
		if dataByte, err = gredis.Get(key); err != nil {
			return
		}
		exist = true
	}
	return
}

func (h *Base) setDataWithKey(key string, data interface{}) error {
	return gredis.Set(key, data, 3600)
}
