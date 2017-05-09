package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func init() {
	control := visitorController{}
	utils.V1.GET("/visitor/tags", control.tags)
}

type visitorController struct {
	baseController
}

func (control visitorController) tags(c *gin.Context) {
	control.initialization(c)
	models := []models.VisitorTag{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)

	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			pipeline = (utils.BsonCreator{}).
				Append(bson.M{"$addFields": bson.M{"name": "$" + control.lang}}).
				Pipeline
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}
