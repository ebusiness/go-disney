package utils

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	//Route - The engine of gin
	Route = gin.Default()
	// version 1 group
	V1 *gin.RouterGroup
	// version 2 group
	V2 *gin.RouterGroup
)

func init() {
	// session && mongo
	Route.Use(middleware.SessionRedisStore, middleware.MongoSession, middleware.CloseNotify)
	// version 1
	V1 = Route.Group("/v1/:lang/:park")
	// version 2
	// V2 = Route.Group("/v2/:lang/:park")

	app := appBase{}
	Route.GET("/", app.index)
	Route.GET("/versions", app.versions)
}

type appBase struct{}

func (app appBase) index(c *gin.Context) {
	c.JSON(http.StatusOK, "disney navigation version 1.0.0")
}

func (app appBase) versions(c *gin.Context) {
	type status struct {
		Version   string `json:"version"`
		Available bool   `json:"available"`
	}
	value := []status{
		status{
			"v1",
			V1 != nil,
		},
		status{
			"v2",
			V2 != nil,
		},
	}
	c.JSON(http.StatusOK, value)
}
