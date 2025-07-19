package config

type Database struct {
	DBHost         string `env:"DB_HOST"`
	DBPort         string `env:"DB_PORT"`
	DBUsername     string `env:"DB_USERNAME"`
	DBPassword     string `env:"DB_PASSWORD"`
	DBDBName       string `env:"DB_DBNAME"`
	IdleConns      string `env:"IDLE_CONNS"`
	MaxConns       string `env:"MAX_CONNS"`
}

func NewDatabase() *Database {
	d := &Database{}
	LoadEnv()
	MarshalEnv(d)
	return d
}
