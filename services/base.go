package services

import (
	"shop/pkg/gredis"
)

func getDataFromRedis(key string) (dataByte []byte, exist bool, err error) {
	if gredis.Exists(key) {
		if dataByte, err = gredis.Get(key); err != nil {
			return
		}
		exist = true
	}
	return
}

func setDataWithKey(key string, data interface{}) error {
	return gredis.Set(key, data, 3600)
}

func setDataWithKeyWithoutExpire(key string, data interface{}) error {
	return gredis.Set(key, data, 0)
}

func deleteCache(key string) (bool, error) {
	return gredis.Delete(key)
}
