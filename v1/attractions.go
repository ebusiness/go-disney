package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	// "log"
)

func init() {
	control := attractionController{}
	utils.V1.GET("/attractions", control.list)
	utils.V1.GET("/attractions/:id", control.detail)
}

type attractionController struct {
	baseController
}

func (control attractionController) commonProject(c *gin.Context, custom bson.M) bson.M {
	if len(control.lang) == 0 {
		control.initialization(c)
	}
	project := bson.M{
		"str_id":           1,
		"main_visual_urls": 1,
		"is_fastpass":      1,
		"thum_url_pc":      1,
		"maps":             1,
		"note":             "$note." + control.lang,
		"introductions":    "$introductions." + control.lang,
		"name":             "$name." + control.lang,
		"area":             1,
	}

	for bsonkey, bsonvalue := range custom {
		project[bsonkey] = bsonvalue
	}
	return bson.M{"$project": project}
}

func (control attractionController) list(c *gin.Context) {
	control.initialization(c)
	conditions := append([]bson.M{},
		bson.M{"$match": bson.M{"park_kind": control.park}},
		control.commonProject(c, nil))
	control.search(c, control.lang, conditions...)
}

func (control attractionController) search(c *gin.Context, lang string, conditions ...bson.M) {
	models := []models.Attraction{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			pipeline = (utils.BsonCreater{}).
				Append(conditions...).
				LookupWithUnwind("areas", "area", "_id", "area", lang).
				// LookupWithUnwind("places", "park", "_id", "park", lang).
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
	control.initialization(c)
	if len(control.id) < 1 || !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	model := models.Attraction{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(model)
	var pipeline []bson.M

	// result := bson.M{}

	utils.SafelyExecutorForGin(c,
		func() {
			basonMatchID := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(control.id)}}
			project := control.commonProject(c, bson.M{"tag_ids": 1, "summary_tag_ids": 1, "summaries": 1})

			pipeline = (utils.BsonCreater{}).
				Append(basonMatchID, project).
				LookupWithUnwind("areas", "area", "_id", "area", control.lang).
				// LookupWithUnwind("places", "park", "_id", "park", control.lang).
				GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags", control.lang).
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
			"_id":  bson.M{"id": "$_id", "type": "$typeid"},
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
	return (utils.BsonCreater{}).Append(bson.M{"$addFields": bson.M{"old": "$$ROOT"}}).
		Append(bson.M{"$unwind": "$summary_tag_ids"}).
		LookupWithUnwind("tagtypes", "summary_tag_ids.typeid", "_id", "type", control.lang).
		LookupWithUnwind("tags", "summary_tag_ids.tagIds", "_id", "tags", control.lang).
		Append(bson.M{"$project": bson.M{"_id": 1, "old": 1, "type": "$type", "tags": "$tags", "typeid": "$summary_tag_ids.typeid"}}).
		Append(groupByRootIDAndTypeID).
		Append(groupByRootID).
		Append(bson.M{"$unwind": "$old.summaries"}).
		Append(bson.M{
			"$group": bson.M{
				"_id": "$_id",
				"old": bson.M{"$first": "$old"},
				"summaries": bson.M{
					"$push": bson.M{
						"body":  "$old.summaries.body.cn",
						"title": "$old.summaries.title.cn",
					},
				},
			},
		}).
		Append(bson.M{"$addFields": bson.M{"old.summary_tags": "$summary_tags", "old.summaries": "$summaries"}}).
		Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$old"}}).
		Pipeline
}
