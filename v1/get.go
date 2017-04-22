package v1

import (
	"fmt"
	"github.com/ebusiness/go-disney/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("get init")
	utils.V1.GET("/test", set)
	utils.V1.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	// route.GET("/get", func(c *gin.Context) {
	// 	session := sessions.Default(c)
	// 	var count int
	// 	v := session.Get("count")
	// 	if v != nil {
	// 		count = v.(int)
	// 	}
	// 	c.JSON(200, gin.H{"count": count})
	// })
}

func set(c *gin.Context) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v != nil {
		count = v.(int)
	}
	c.JSON(200, gin.H{"count": count})
}
