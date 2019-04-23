package redigoSentinelstore

import (
	"errors"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"time"
)

var MySentinel *redis.Pool
var MyRedisConf *RedisConf

func RedisInit(c *RedisConf) *redis.Pool {
	sntnl := &sentinel.Sentinel{
		Addrs:      c.Sentinels,
		MasterName: c.MasterName,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := time.Millisecond * time.Duration(c.TimeOut)
			c, err := redis.Dial("tcp", addr,
				redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout),
				redis.DialWriteTimeout(timeout))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	MySentinel = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: time.Millisecond * time.Duration(c.TimeOut),
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr,
				redis.DialConnectTimeout(time.Millisecond*time.Duration(c.TimeOut)),
				redis.DialReadTimeout(time.Millisecond*time.Duration(c.TimeOut)),
				redis.DialWriteTimeout(time.Millisecond*time.Duration(c.TimeOut)))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(rc redis.Conn, t time.Time) error {
			if !Auth(rc, c.Password) {
				return errors.New("Auth failed")
			}

			if !sentinel.TestRole(rc, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}
	return MySentinel
}

func Auth(c redis.Conn, password string) bool {
	result, err := redis.Bool(c.Do("AUTH", password))
	if err != nil || !result {
		return false
	}
	return true
}
