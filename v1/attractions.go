package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
)

func init() {
	control := attractionController{}
	utils.V1.GET("/attractions", control.list)
	utils.V1.GET("/attractions/:id", control.detail)
}

type attractionController struct{}

func (control attractionController) list(c *gin.Context) {
	models := []models.Attraction{}

	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)

	pipeline := (utils.BsonCreater{}).
		LookupWithUnwind("areas", "area", "_id", "area").
		LookupWithUnwind("places", "park", "_id", "park").
		// GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags").
		Pipeline

	if c.IsAborted() {
		return
	}
	collection.Pipe(pipeline).All(&models)

	if c.IsAborted() {
		return
	}
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Println("something wrong", err)
	// 	}
	// }()
	c.JSON(http.StatusOK, models)
}

func (control attractionController) detail(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 1 || !bson.IsObjectIdHex(id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	log.Println(id)

	model := models.Attraction{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(model)

	basonMatchID := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}
	tempRoot := bson.M{"$addFields": bson.M{"old": "$$ROOT"}}

	pipeline := (utils.BsonCreater{}).
		Append(basonMatchID).
		LookupWithUnwind("areas", "area", "_id", "area").
		LookupWithUnwind("places", "park", "_id", "park").
		GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags").
		Append(tempRoot).
		Append(control.lookupSummaryTags()...).
		Pipeline
	if c.IsAborted() {
		return
	}
	collection.Pipe(pipeline).One(&model)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, model)
}

func (control attractionController) lookupSummaryTags() []bson.M {
	groupByRootIDAndTypeID := bson.M{
		"$group": bson.M{
			"_id":  bson.M{"id": "$_id", "type": "$type._id"},
			"old":  bson.M{"$first": "$old"},
			"type": bson.M{"$first": "$type"},
			"tags": bson.M{"$push": "$tags"},
		},
	}
	groupByRootID := bson.M{
		"$group": bson.M{
			"_id": "$_id.id",
			"old": bson.M{"$first": "$old"},
			"summary_tags": bson.M{
				"$push": bson.M{"type": "$type", "tags": "$tags"},
			},
		},
	}

	creater := utils.BsonCreater{}
	return creater.Append(bson.M{"$unwind": "$summary_tag_ids"}).
		LookupWithUnwind("tagtypes", "summary_tag_ids.typeid", "_id", "type").
		LookupWithUnwind("tags", "summary_tag_ids.tagIds", "_id", "tags").
		Append(bson.M{"$project": bson.M{"_id": 1, "old": 1, "type": "$type", "tags": "$tags"}}).
		Append(groupByRootIDAndTypeID).
		Append(groupByRootID).
		Append(bson.M{"$addFields": bson.M{"old.summary_tags": "$summary_tags"}}).
		Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$old"}}).
		Pipeline
}
