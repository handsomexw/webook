//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/webook?charset=utf8&parseTime=True&loc=Local",
	},
	Redisaddr: RedisConfig{
		Addr: "localhost:16379",
	},
}
