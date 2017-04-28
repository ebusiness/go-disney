package utils

import (
	"github.com/gin-gonic/gin"

	"github.com/ebusiness/go-disney/middleware"
)

var (
	//Route - The engine of gin
	Route = gin.Default()
	//V1 - Route version 1
	V1 *gin.RouterGroup //= Route.Group("/v1")
)

func init() {
	// session && mongo
	Route.Use(middleware.SessionRedisStore, middleware.MongoSession, middleware.CloseNotify)
	// version 1
	V1 = Route.Group("/v1")

	// V1.Use(middleware.SessionRedisStore())
}
