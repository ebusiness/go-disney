package utils

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
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
	V1 = Route.Group("/v1/:lang/:park")

	app := appBase{}
	Route.GET("/", app.index)
	Route.GET("/versions", app.versions)
}

type appBase struct{}

func (app appBase) index(c *gin.Context) {
	c.JSON(http.StatusOK, "disney navigation version 1.0.0")
}

func (app appBase) versions(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]bool{
		"V1": V1 != nil,
	})
}
