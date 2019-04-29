package redigoSentinel

import (
	"errors"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisConf struct {
	Network    string
	Password   string
	DB         int
	TimeOut    int
	Pool       int
	MasterName string
	Sentinels  []string
}

var MySentinel *redis.Pool
var MyRedisConf *RedisConf

func RedisInit(c *RedisConf) (*redis.Pool, error) {
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
			var rc redis.Conn
			if c.Password != "" {
				rc, err = redis.Dial("tcp", masterAddr,
					redis.DialConnectTimeout(time.Millisecond*time.Duration(c.TimeOut)),
					redis.DialReadTimeout(time.Millisecond*time.Duration(c.TimeOut)),
					redis.DialWriteTimeout(time.Millisecond*time.Duration(c.TimeOut)),
					redis.DialPassword(c.Password), redis.DialDatabase(c.DB))
			} else {
				rc, err = redis.Dial("tcp", masterAddr,
					redis.DialConnectTimeout(time.Millisecond*time.Duration(c.TimeOut)),
					redis.DialReadTimeout(time.Millisecond*time.Duration(c.TimeOut)),
					redis.DialWriteTimeout(time.Millisecond*time.Duration(c.TimeOut)), redis.DialDatabase(c.DB))
			}

			if err != nil {
				return nil, err
			}
			return rc, nil
		},
		TestOnBorrow: func(rc redis.Conn, t time.Time) error {
			if !sentinel.TestRole(rc, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}
	conn := MySentinel.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("PING"))
	if err != nil || r != "PONG" {
		return nil, errors.New("PING check failed")
	}
	return MySentinel, nil
}
