package hjsentinel

import (
	"errors"
	"fmt"
	"github.com/mediocregopher/radix"
	"log"
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

func sentinelConnFunc(network, addr string) (radix.Conn, error) {
	c, err := radix.Dial(network, addr, radix.DialTimeout(time.Millisecond*time.Duration(MyRedisConf.TimeOut)))
	if err != nil {
		log.Println("sentinel conn dial 报错！！！err:", err)
		return c, err
	}
	var ping_result string
	err = c.Do(radix.Cmd(&ping_result, "ping"))
	if err != nil {
		log.Println("sentinel conn do ping 报错！！！err:", err)
		return c, err
	}
	if ping_result != "PONG" {
		log.Println("sentinel conn do ping 没有受到PONG！！！ping_result:", ping_result)
		return c, errors.New(fmt.Sprintln("sentinel conn do ping 没有受到PONG！！！ping_result:", ping_result))
	}
	log.Println("连接sentine 成功！！！")
	return c, nil
}

// redis 连接函数
func redisConnFunc(network, addr string) (radix.Conn, error) {
	c, err := radix.Dial(network, addr, radix.DialAuthPass(MyRedisConf.Password), radix.DialSelectDB(MyRedisConf.DB),
		radix.DialTimeout(time.Millisecond*time.Duration(MyRedisConf.TimeOut)))
	if err != nil {
		log.Println("redis conn dial 报错！！！err:", err)
		return c, err
	}
	var ping_result string
	err = c.Do(radix.Cmd(&ping_result, "ping"))
	if err != nil {
		log.Println("redis conn do ping 报错！！！err:", err)
		return c, err
	}
	if ping_result != "PONG" {
		log.Println("redis conn do ping 没有受到PONG！！！ping_result:", ping_result)
		return c, errors.New(fmt.Sprintln("redis conn do ping 没有受到PONG！！！ping_result:", ping_result))
	}
	log.Println("连接 redis 成功！！！")
	return c, nil
}

func redisPoolFunc(network, addr string) (radix.Client, error) {

	p, err := radix.NewPool(network, addr, MyRedisConf.Pool, radix.PoolConnFunc(redisConnFunc), radix.PoolOnFullClose())
	if err != nil {
		log.Println("redis pool dial 报错！！！errr:", err)
		return p, err
	}
	return p, nil
}

var Sentinel *radix.Sentinel
var MyRedisConf *RedisConf

func RedisInit(c *RedisConf) (*radix.Sentinel, error) {
	MyRedisConf = c
	scf := radix.SentinelConnFunc(sentinelConnFunc)
	spf := radix.SentinelPoolFunc(redisPoolFunc)
	sentinel, err := radix.NewSentinel(MyRedisConf.MasterName, MyRedisConf.Sentinels, scf, spf)
	if err != nil {
		log.Println("创建 new sentinel 报错！！！errr:", err)
		return nil, err
	}
	Sentinel = sentinel
	log.Println("Redis 初始化成功！")
	return sentinel, nil
}
