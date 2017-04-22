package v1

import (
	"github.com/ebusiness/go-disney/utils"
	"github.com/gin-gonic/gin"
)

//Regist - regist all controllers of version 1
// just touch Regist(), it will be auto load all `init` function of this package's files
func Regist() {
	utils.V1.GET("/", index)
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{"hello": "world"})
}
