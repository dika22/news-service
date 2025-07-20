package config

type Config struct {
	AppPort                   string `env:"APP_PORT"`
	RateLimitThreshold        string `env:"RATE_LIMIT_THRESHOLD"`
	DebugHTTP                 string `env:"DEBUG_HTTP"`
	MessageBrokerURL          string `env:"MESSAGE_BROKER_URL"`
	MessageBrokerUsername     string `env:"MESSAGE_BROKER_USERNAME"`
	MessageBrokerPassword     string `env:"MESSAGE_BROKER_PASSWORD"`
	MessageBrokerPort         string `env:"MESSAGE_BROKER_PORT"`
	MessageBrokerQueue        string `env:"MESSAGE_BROKER_QUEUE"`
	ElasticsearchURL          string `env:"ELASTICSEARCH_URL"`
	ElasticsearchUsername     string `env:"ELASTICSEARCH_USERNAME"`
	ElasticsearchPassword     string `env:"ELASTICSEARCH_PASSWORD"`
	ArticleQueue              string `env:"ARTICLE_QUEUE"`
	ArticleIndex              string `env:"ARTICLE_INDEX"`
	RateLimitMaxRequest       string `env:"RATE_LIMIT_MAX_REQUEST"`
	RateLimitInterval         string `env:"RATE_LIMIT_INTERVAL"`
	RateLimitJitter           string `env:"RATE_LIMIT_JITTER"`
}

func NewConfig() *Config {
	c := &Config{}
	LoadEnv()
	MarshalEnv(c)
	if c.RateLimitThreshold == "" {
		c.RateLimitThreshold = "1000"
	}
	return c
}
