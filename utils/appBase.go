package utils

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	//Route - The engine of gin
	Route = gin.Default()
	// V1 version 1 group
	V1 *gin.RouterGroup
	// V2 version 2 group
	V2 *gin.RouterGroup
)

func init() {
	// session && mongo
	Route.Use(middleware.SessionRedisStore, middleware.MongoSession, middleware.CloseNotify)

	app := appBase{}
	Route.GET("/", app.index)
	Route.NoRoute(app.notFound)

	Route.GET("/versions", app.versions)

	// version 1
	V1 = Route.Group("/v1/:lang/:park")
	// {
	// 	Route.GET("/v1/:lang", app.index)
	// 	V1.GET("/", app.index)
	// }c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	// version 2
	// V2 = Route.Group("/v2/:lang/:park")
}

type appBase struct{}

func (app appBase) index(c *gin.Context) {
	c.JSON(http.StatusOK, "disney navigation version 1.0.0")
}

func (app appBase) notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, []gin.H{
		app.routers("v1", V1),
	})
}

func (app appBase) routers(version string, group *gin.RouterGroup) gin.H {
	if group == nil {
		return gin.H{}
	}
	type route struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	}
	basePath := group.BasePath()
	routes := []route{}
	for _, item := range Route.Routes() {
		if !strings.Contains(item.Path, basePath) {
			continue
		}
		routes = append(routes, route{
			item.Method,
			strings.Replace(item.Path, basePath, "", -1),
		})
	}
	return gin.H{
		"version":  version,
		"basePath": basePath,
		"routers":  routes,
	}
}

func (app appBase) versions(c *gin.Context) {
	type status struct {
		Version   string `json:"version"`
		Available bool   `json:"available"`
	}
	value := []status{
		{
			"v1",
			V1 != nil,
		},
		{
			"v2",
			V2 != nil,
		},
	}
	c.JSON(http.StatusOK, value)
}
