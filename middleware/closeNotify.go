package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//CloseNotify - HTTP connection closed
func CloseNotify(c *gin.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorln("Handler finished without response body", err)
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
