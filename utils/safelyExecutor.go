package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// SafelyExecutorForGin - execute function safely for gin
func SafelyExecutorForGin(c *gin.Context, tasks ...func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("something wrong", err) // stream closed
			c.AbortWithStatus(http.StatusNotAcceptable)
		}
	}()

	for _, task := range tasks {
		if c.IsAborted() {
			return
		}
		task()
	}
}
