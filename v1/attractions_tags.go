package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	// "github.com/ebusiness/go-disney/v1/algorithms"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	// "log"
	// log "github.com/sirupsen/logrus"
	"net/http"
	// "time"
)

func init() {
	control := attractionController{}
	utils.V1.GET("/attraction/tags", control.tags)
	utils.V1.GET("/attraction/tags/:id", control.attractionsOfTag)
}

func (control attractionController) tags(c *gin.Context) {
	control.initialization(c)
	models := []models.AttractionTag{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			pipeline = (utils.BsonCreator{}).
				Append(bson.M{"$project": bson.M{"_id": 1, "icon": 1, "color": 1, "name": "$name." + control.lang}}).
				Pipeline
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}

func (control attractionController) attractionsOfTag(c *gin.Context) {
	control.initialization(c)
	if len(control.id) < 1 || !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	models := []models.AttractionTag{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			pipeline = (utils.BsonCreator{}).
				Append(bson.M{"$project": bson.M{"_id": 1, "icon": 1, "color": 1, "name": "$name." + control.lang}}).
				Pipeline
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}
