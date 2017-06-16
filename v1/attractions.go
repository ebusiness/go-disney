package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
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
		"area":             1,
		"is_lottery":       1,
		"is_must_book":     1,
		"category":         1,
		"is_available":     1,
		"realtime":         1,
		"note":             "$note." + control.lang,
		"introductions":    "$introductions." + control.lang,
		"name":             "$name." + control.lang,
	}

	for bsonkey, bsonvalue := range custom {
		project[bsonkey] = bsonvalue
	}
	return bson.M{"$project": project}
}

func (control attractionController) list(c *gin.Context) {
	///////////////////////
	async := utils.Executor(c)
	log.Println(async)
	///////////////////////
	control.initialization(c)
	getPipeline := func(param interface{}) (interface{}, error) {
		pipeline := (utils.BsonCreator{}).
			Append(bson.M{"$match": bson.M{"park_kind": control.park, "name." + control.lang: bson.M{"$ne": ""}}}).
			Append(control.commonProject(c, nil)).
			LookupWithUnwind("areas", "area", "_id", "area", control.lang).
			// LookupWithUnwind("places", "park", "_id", "park", lang).
			// GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags").
			Pipeline

		sortQuery := c.Query("sort")
		if sortQuery == "hot" {
			pipeline = (utils.BsonCreator{}).
				Append(pipeline...).
				LookupWithUnwind("attractions_hot", "str_id", "str_id", "hot", "").
				Append(bson.M{"$addFields": bson.M{"index_hot": "$hot.hot"}}).
				Append(bson.M{"$sort": bson.M{"hot.hot": -1}}).
				Pipeline
		}
		return pipeline, nil
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.([]bson.M)
		models := []models.Attraction{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollection(models)

		// res := []bson.M{}
		// log.Println("test")
		err := collection.Pipe(pipeline).All(&models)
		return models, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
}

func (control attractionController) detail(c *gin.Context) {
	control.initialization(c)
	if len(control.id) < 1 { //|| !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	getPipeline := func(param interface{}) (interface{}, error) {
		basonMatchID := bson.M{"$match": bson.M{"str_id": control.id}} //bson.ObjectIdHex(control.id)}}
		project := control.commonProject(c, bson.M{"tag_ids": 1, "limited": 1, "youtube_url": 1, "summary_tag_ids": 1, "summaries": 1})

		return (utils.BsonCreator{}).
			Append(basonMatchID, project).
			LookupWithUnwind("areas", "area", "_id", "area", control.lang).
			GraphLookup("limiteds", "$limited", "limited", "_id", "limited", control.lang).
			GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags", control.lang).
			Append(control.lookupSummaryTags()...).
			LookupWithUnwind("attractions_hot", "str_id", "str_id", "hot", "").
			Append(bson.M{"$addFields": bson.M{"index_hot": "$hot.hot"}}).
			Pipeline, nil
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.([]bson.M)
		model := models.Attraction{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollection(model)
		err := collection.Pipe(pipeline).One(&model)
		return model, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
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
	return (utils.BsonCreator{}).Append(bson.M{"$addFields": bson.M{"old": "$$ROOT"}}).
		Append(bson.M{"$unwind": "$summary_tag_ids"}).
		LookupWithUnwind("tagtypes", "summary_tag_ids.typeid", "_id", "type", control.lang).
		LookupWithUnwind("tags", "summary_tag_ids.tagIds", "_id", "tags", control.lang).
		Append(bson.M{"$project": bson.M{"_id": 1, "old": 1, "type": "$type", "tags": "$tags", "typeid": "$summary_tag_ids.typeid"}}).
		Append(groupByRootIDAndTypeID).
		Append(groupByRootID).
		Append(bson.M{"$unwind": bson.M{"path": "$old.summaries", "preserveNullAndEmptyArrays": true}}).
		Append(bson.M{
			"$group": bson.M{
				"_id":          "$_id",
				"old":          bson.M{"$first": "$old"},
				"summary_tags": bson.M{"$first": "$summary_tags"},
				"summaries": bson.M{
					"$addToSet": bson.M{
						"body":  "$old.summaries.body." + control.lang,
						"title": "$old.summaries.title." + control.lang,
					},
				},
			},
		}).
		Append(bson.M{"$addFields": bson.M{"old.summary_tags": "$summary_tags", "old.summaries": bson.M{
			"$cond": bson.M{
				"if": bson.M{"$eq": []interface{}{
					0, bson.M{
						"$size": "$summaries.body",
					},
				},
				},
				"then": nil,
				"else": "$summaries",
			},
		}}}).
		Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$old"}}).
		Pipeline
}
