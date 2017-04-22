package config

const (
	// A default server only if os.Getenv("REDISSERVER") is not set
	RedisDefaultServer = "localhost"
	// A default port only if os.Getenv("REDISPORT") is not set
	RedisDefaultPort = "6379"
)

const (
	// Session key
	SessionKey = "go-disney"
	// session secret key [Aes-256]
	SecretAes = "yVHlew1jDlZpJ/zSbJ8JPjIc2dBeoLny"
	// The max age of session
	MaxAge = 86400 * 30
)
