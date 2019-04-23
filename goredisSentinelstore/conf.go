package goredisSentinelstore

type RedisConf struct {
	Password   string
	DB         int
	TimeOut    int
	Pool       int
	MasterName string
	Sentinels  []string
}
