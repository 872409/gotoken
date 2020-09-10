package gotoken



type GoTokenConfig struct {
	RedisHost string
	RedisPwd  string
	RedisDB   int
	TokenVersion string
	Secret       string
	ExpireHour   int
}
