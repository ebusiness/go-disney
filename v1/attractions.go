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
	control := attractionController{}
	utils.V1.GET("/land/attractions", control.landlist)
	utils.V1.GET("/sea/attractions", control.sealist)
	utils.V1.GET("/attractions", control.list)
	utils.V1.GET("/attractions/:id", control.detail)
}

type attractionController struct{}

func (control attractionController) landlist(c *gin.Context) {
	park := bson.M{"$match": bson.M{"park_kind": "1"}}
	control.search(c, park)
}

func (control attractionController) sealist(c *gin.Context) {
	park := bson.M{"$match": bson.M{"park_kind": "2"}}
	control.search(c, park)
}

func (control attractionController) list(c *gin.Context) {
	// park := bson.M{"$match": bson.M{"park_kind": bson.M{"$in": []string{"1", "2"}}}}
	control.search(c)
}

func (control attractionController) search(c *gin.Context, bsons ...bson.M) {
	models := []models.Attraction{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			pipeline = (utils.BsonCreater{}).
				Append(bsons...).
				LookupWithUnwind("areas", "area", "_id", "area").
				LookupWithUnwind("places", "park", "_id", "park").
				// GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags").
				Pipeline
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}

func (control attractionController) detail(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 1 || !bson.IsObjectIdHex(id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	model := models.Attraction{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(model)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			basonMatchID := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(id)}}
			tempRoot := bson.M{"$addFields": bson.M{"old": "$$ROOT"}}
			pipeline = (utils.BsonCreater{}).
				Append(basonMatchID).
				LookupWithUnwind("areas", "area", "_id", "area").
				LookupWithUnwind("places", "park", "_id", "park").
				GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags").
				Append(tempRoot).
				Append(control.lookupSummaryTags()...).
				Pipeline
		},
		func() {
			collection.Pipe(pipeline).One(&model)
		},
		func() {
			c.JSON(http.StatusOK, model)
		})
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
