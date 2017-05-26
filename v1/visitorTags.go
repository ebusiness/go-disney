package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
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
	getPipeline := func(param interface{}) (interface{}, error) {
		return (utils.BsonCreator{}).
			Append(bson.M{"$addFields": bson.M{"name": "$" + control.lang}}).
			Pipeline, nil
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.([]bson.M)
		models := []models.VisitorTag{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollection(models)
		err := collection.Pipe(pipeline).All(&models)
		return models, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
}
