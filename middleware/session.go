package middleware

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/ebusiness/go-disney/config"
)

//SessionRedisStore - the middleware of session store by Redis
var SessionRedisStore = func() gin.HandlerFunc {
	//store := sessions.NewCookieStore([]byte("secret"))
	store, err := sessions.NewRedisStore(10, "tcp", getRedisServerAddress(), "", []byte(config.SecretAes))
	if err != nil {
		log.Fatalln("cannot connect to the redis server", err)
	}
	store.Options(sessions.Options{
		MaxAge: config.MaxAge,
	})
	return sessions.Sessions(config.SessionKey, store)
}()

func getRedisServerAddress() string {
	redisServer := os.Getenv("REDISSERVER")
	if len(redisServer) < 1 {
		redisServer = config.RedisDefaultServer
	}
	redisPort := os.Getenv("REDISPORT")
	if len(redisPort) < 1 {
		redisPort = config.RedisDefaultPort
	}
	return redisServer + ":" + redisPort
}
