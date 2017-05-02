package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

//Regist - regist all controllers of version 1
// just touch Regist(), it will be auto load all `init` function of this package's files
func init() {
	utils.V1.GET("/visitor/tags", visitorIndex)
}

func visitorIndex(c *gin.Context) {
	models := []models.VisitorTag{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)

	utils.SafelyExecutorForGin(c,
		func() {
			collection.Find(bson.M{}).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}
