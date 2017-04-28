package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//CloseNotify - HTTP connection closed
func CloseNotify(c *gin.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Handler finished without response body", err)
				// log.Println("Error:", err)
			}
		}()

		<-c.Writer.CloseNotify()
		if c.Writer.Written() {
			return
		}

		if c.IsAborted() {
			return
		}
		c.AbortWithStatus(http.StatusNotAcceptable) //406
	}()
	c.Next()
}
