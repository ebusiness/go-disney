package v1

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
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
	collection.Find(bson.M{}).All(&models)

	c.JSON(200, models)
}
