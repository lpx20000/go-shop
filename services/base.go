package services

import (
	"shop/pkg/gredis"
)

func get(key string) (dataByte []byte, exist bool, err error) {
	if gredis.Exists(key) {
		if dataByte, err = gredis.Get(key); err != nil {
			return
		}
		exist = true
	}
	return
}

func set(key string, data interface{}) error {
	return gredis.Set(key, data, 3600)
}

func setWithoutExpire(key string, data interface{}) error {
	return gredis.Set(key, data, 0)
}

func deleteCache(key string) (bool, error) {
	return gredis.Delete(key)
}

func hget(key, field string) (reply []byte, err error) {
	reply, err = gredis.Hget(key, field)
	return
}

func hset(key, field string, data interface{}) error {
	return gredis.Hset(key, field, data)
}

func hdel(key, field string) (bool, error) {
	return gredis.Delete(key)
}
