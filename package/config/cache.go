package config

type Cache struct {
	RedisHost       string `env:"REDIS_HOST"`
	RedisPort       string `env:"REDIS_PORT"`
	RedisPrefix     string `env:"CACHE_PREFIX"`
	WorkerRedisHost string `env:"WORKER_REDIS_HOST"`
	WorkerRedisPort string `env:"WORKER_REDIS_PORT"`
	LRUSize         string `env:"LRU_SIZE"`
	MaxIdleConn     string `env:"CACHE_MAX_IDLE_CONN"`
	MinIdleConn     string `env:"CACHE_MIN_IDLE_CONN"`
	PoolSize        string `env:"CACHE_POOL_SIZE"`
}

func NewCache() *Cache {
	c := &Cache{}
	LoadEnv()
	MarshalEnv(c)
	return c
}
