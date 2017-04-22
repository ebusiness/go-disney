package config

const (
	//RedisDefaultServer - A default server only if os.Getenv("REDISSERVER") is not set
	RedisDefaultServer = "localhost"
	//RedisDefaultPort - A default port only if os.Getenv("REDISPORT") is not set
	RedisDefaultPort = "6379"
)

const (
	//SessionKey - Session key
	SessionKey = "go-disney"
	//SecretAes - session secret key [Aes-256]
	SecretAes = "yVHlew1jDlZpJ/zSbJ8JPjIc2dBeoLny"
	//MaxAge - The max age of session
	MaxAge = 86400 * 30
)
