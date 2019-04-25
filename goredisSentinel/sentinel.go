package goredisSentinel

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

type RedisConf struct {
	Password   string
	DB         int
	TimeOut    int
	Pool       int
	MasterName string
	Sentinels  []string
}

var MySentinel *redis.Client
var MyRedisConf *RedisConf

func RedisInit(rc *RedisConf) (*redis.Client, error) {
	timeout := time.Millisecond * time.Duration(rc.TimeOut)
	fo := &redis.FailoverOptions{MasterName: rc.MasterName, SentinelAddrs: rc.Sentinels,
		Password: rc.Password, DB: rc.DB, DialTimeout: timeout, ReadTimeout: timeout,
		WriteTimeout: timeout, PoolTimeout: timeout, IdleTimeout: timeout * 1000, PoolSize: rc.Pool}
	c := redis.NewFailoverClient(fo)
	r, err := c.Ping().Result()
	if r != "PONG" || err != nil {
		return nil, errors.New("PING check failed")
	}
	return c, nil
}
