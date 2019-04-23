package goredisSentinelstore

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

var MySentinel *redis.Client
var MyRedisConf *RedisConf

func RedisInit(rc *RedisConf) (*redis.Client, error) {
	timeout := time.Millisecond * time.Duration(rc.TimeOut)
	fo := &redis.FailoverOptions{MasterName: rc.MasterName, SentinelAddrs: rc.Sentinels,
		Password: rc.Password, DB: rc.DB, DialTimeout: timeout, ReadTimeout: timeout,
		WriteTimeout: timeout, PoolTimeout: timeout, IdleTimeout: timeout, PoolSize: rc.Pool}
	c := redis.NewFailoverClient(fo)
	r, err := c.Ping().Result()
	if r != "PONG" || err != nil {
		return nil, errors.New("PING check failed")
	}
	MySentinel = c
	return c, nil
}
