//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:13307)/webook?charset=utf8&parseTime=True&loc=Local",
	},
	Redisaddr: RedisConfig{
		Addr: "webook-redis:11379",
	},
}
