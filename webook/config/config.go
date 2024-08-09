package config

type config struct {
	DB        DBConfig
	Redisaddr RedisConfig
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}
