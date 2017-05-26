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

	getPipeline := func(param interface{}) (interface{}, error) {
		return (utils.BsonCreator{}).
			Append(bson.M{"$project": bson.M{"_id": 1, "plan_creation": 1, "name": "$name." + control.lang}}).
			Pipeline, nil
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.([]bson.M)
		models := []models.AttractionTag{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollection(models)
		err := collection.Pipe(pipeline).All(&models)
		return models, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
}

func (control attractionController) attractionsOfTag(c *gin.Context) {
	control.initialization(c)
	if len(control.id) < 1 || !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	getConditions := func(param interface{}) (interface{}, error) {
		return append([]bson.M{},
			bson.M{"$match": bson.M{"park_kind": control.park, "name." + control.lang: bson.M{"$ne": ""}}},
			control.commonProject(c, bson.M{"tempid": 1})), nil
	}
	getPipeline := func(param interface{}) (interface{}, error) {
		conditions := param.([]bson.M)
		return (utils.BsonCreator{}).
			Append(bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(control.id)}}).
			Append(bson.M{"$unwind": "$attractions"}).
			Append(bson.M{
				"$project": bson.M{
					"str_id": "$attractions.str_id",
					"rank":   bson.M{"$cmp": []string{"$attractions.ranking", "$plan_creation.threshold"}},
				},
			}).
			Append(bson.M{"$match": bson.M{"rank": bson.M{"$gt": 0}}}).
			LookupWithUnwind("attractions", "str_id", "str_id", "attraction", "").
			Append(bson.M{"$addFields": bson.M{"attraction.tempid": "$_id"}}).
			Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$attraction"}}).
			Append(conditions...).
			LookupWithUnwind("areas", "area", "_id", "area", control.lang).
			Append(bson.M{
				"$project": bson.M{
					"_id":        "$tempid",
					"attraction": "$$ROOT",
				},
			}).
			Append(bson.M{
				"$group": bson.M{
					"_id":         "$_id",
					"attractions": bson.M{"$push": "$attraction"},
				},
			}).
			Pipeline, nil
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.([]bson.M)
		model := models.AttractionTag{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollection(model)
		err := collection.Pipe(pipeline).One(&model)
		if len(model.Attractions) == 0 {
			return nil, err
		}
		return model, err
	}
	utils.Executor(c).Waterfall(getConditions, getPipeline, exec)
}
