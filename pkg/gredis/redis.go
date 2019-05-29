package gredis

import (
	"encoding/json"
	"shop/pkg/logging"
	"shop/pkg/setting"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

func SetUp() {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				logging.LogFatal(err.Error())
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					_ = c.Close()
					logging.LogFatal(err.Error())
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				logging.LogFatal(err.Error())
				return err
			}
			return nil
		},
	}
}

func Set(key string, data interface{}, time int) (err error) {
	var (
		value []byte
		conn  redis.Conn
	)

	if setting.RedisSetting.RedisOpen {
		conn = RedisConn.Get()
		defer conn.Close()

		if value, err = json.Marshal(data); err != nil {
			return err
		}

		if _, err = conn.Do("SET", key, value); err != nil {
			return err
		}

		if _, err = conn.Do("EXPIRE", key, time); err != nil {
			return err
		}
	}
	return nil
}

func Exists(key string) (exist bool) {
	var (
		conn redis.Conn
		err  error
	)

	if setting.RedisSetting.RedisOpen {
		conn = RedisConn.Get()
		defer conn.Close()
		if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
			return false
		}
	}
	return exist
}

func Get(key string) (reply []byte, err error) {
	var (
		conn redis.Conn
	)

	if setting.RedisSetting.RedisOpen {
		conn = RedisConn.Get()
		defer conn.Close()

		if reply, err = redis.Bytes(conn.Do("GET", key)); err != nil {
			return nil, err
		}
	}
	return reply, nil
}

func Delete(key string) (bool, error) {
	if setting.RedisSetting.RedisOpen {
		conn := RedisConn.Get()
		defer conn.Close()
		return redis.Bool(conn.Do("DEL", key))
	}
	return false, nil
}
func LikeDeletes(key string) error {
	if setting.RedisSetting.RedisOpen {
		conn := RedisConn.Get()
		defer conn.Close()

		keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
		if err != nil {
			return err
		}

		for _, key := range keys {
			_, err = Delete(key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
