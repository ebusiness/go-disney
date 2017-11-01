package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

func init() {
	control := showController{}
	utils.V1.GET("/show/:date", control.list)
	utils.V1.GET("/show/:date/:id", control.detail)
}

type showController struct {
	baseController
}

func (control showController) commonProject(c *gin.Context, custom bson.M) bson.M {
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

func (control showController) list(c *gin.Context) {

	dateString := c.Param("date")
	if len(dateString) < 0 {
		c.AbortWithStatus(http.StatusNotAcceptable) //406
		return
	}
	t, err := time.Parse("2006-01-02 -0700", dateString+" +0900")
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable) //406
		return
	}

	///////////////////////
	log.Println(t)
	///////////////////////
	control.initialization(c)

	getPipeline := func(param interface{}) (interface{}, error) {
		bsonCreator := (utils.BsonCreator{}).
			Append(bson.M{"$match": bson.M{
				"park_kind":            control.park,
				"name." + control.lang: bson.M{"$ne": ""},
				"category":             "show",
			}}).
			// Append(bson.M{"$match": bson.M{"str_id": "wb_penny"}}).
			Append(control.commonProject(c, nil)).
			LookupWithUnwind("areas", "area", "_id", "area", control.lang)
			// LookupWithUnwind("places", "park", "_id", "park", lang).

		pipeline := control.getPipelineForShowSchedules(t, bsonCreator).Pipeline
		pipeline = control.joinLocation(pipeline)

		///////////////////////
		// log.Println(pipeline)
		///////////////////////
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

		// err := collection.Pipe(pipeline).All(&models)
		// return models, err

		res := []bson.M{}
		// log.Println("test")
		err := collection.Pipe(pipeline).All(&res)
		return res, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
}

func (control showController) getPipelineForShowSchedules(t time.Time, bsonCreator utils.BsonCreator) utils.BsonCreator {
	return bsonCreator.Append(bson.M{"$addFields": bson.M{"old": "$$ROOT"}}).
		Lookup("schedule_daily", "str_id", "str_id", "schedule").
		Append(bson.M{"$addFields": bson.M{"schedules": bson.M{
			"$cond": bson.M{
				"if": bson.M{"$eq": []interface{}{
					0, bson.M{
						"$size": "$schedule",
					},
				},
				},
				"then": []bson.M{{"str_id": "$str_id"}},
				"else": "$schedule",
			},
		},
			"schedule": nil,
		}}).
		Append(
			bson.M{
				"$unwind": "$schedules",
			},
		).
		Append(bson.M{"$addFields": bson.M{
			"schedules": bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$and": []bson.M{
							{"$lt": []interface{}{
								"$schedules.startTime", t.AddDate(0, 0, 1),
							}},
							{"$gt": []interface{}{
								"$schedules.startTime", t,
							}},
						},
					},
					"then": "$schedules",
					"else": "$unset",
				},
			},
			"schedule": nil,
		}}).
		Append(bson.M{
			"$group": bson.M{
				"_id": "$str_id",
				"old": bson.M{"$first": "$old"},
				"schedules": bson.M{
					"$push": "$schedules",
				},
			},
		}).
		Append(bson.M{"$addFields": bson.M{"old.schedules": "$schedules"}}).
		Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$old"}})
}

func (control showController) detail(c *gin.Context) {
	control.initialization(c)
	if len(control.id) < 1 { //|| !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	dateString := c.Param("date")
	if len(dateString) < 0 {
		c.AbortWithStatus(http.StatusNotAcceptable) //406
		return
	}
	t, err := time.Parse("2006-01-02 -0700", dateString+" +0900")
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable) //406
		return
	}
	log.Println(t)

	getPipeline := func(param interface{}) (interface{}, error) {
		basonMatchID := bson.M{"$match": bson.M{"str_id": control.id, "category": "show"}} //bson.ObjectIdHex(control.id)}}
		project := control.commonProject(c, bson.M{"tag_ids": 1, "limited": 1, "youtube_url": 1, "summary_tag_ids": 1, "summaries": 1})

		bsonCreator := (utils.BsonCreator{}).
			Append(basonMatchID, project).
			LookupWithUnwind("areas", "area", "_id", "area", control.lang).
			GraphLookup("limiteds", "$limited", "limited", "_id", "limited", control.lang).
			GraphLookup("tags", "$tag_ids", "tag_ids", "_id", "tags", control.lang).
			Append(control.lookupSummaryTags()...).
			LookupWithUnwind("attractions_hot", "str_id", "str_id", "hot", "").
			Append(bson.M{"$addFields": bson.M{"index_hot": "$hot.hot"}})

		pipeline := control.getPipelineForShowSchedules(t, bsonCreator).Pipeline
		pipeline = control.joinLocation(pipeline)

		return pipeline, nil
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

func (control showController) lookupSummaryTags() []bson.M {
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



func (control showController) joinLocation(pipeline []bson.M) []bson.M {
	log.Println("attraction_location")
	return (utils.BsonCreator{}).
		Append(pipeline...).
		Lookup("attraction_location", "str_id", "str_id", "location").
		Append(bson.M{"$addFields": bson.M{
			"location": bson.M{
				"$cond": bson.M{
					"if": bson.M{"$eq": []interface{}{
						0, bson.M{
							"$size": "$location",
						},
					},
					},
					"then": []bson.M{{"str_id": "$str_id"}},
					"else": "$location",
				},
			},
		}}).
		Append(
			bson.M{
				"$unwind": "$location",
			},
		).
		Append(bson.M{"$addFields": bson.M{
			"coordinates" : "$location.coordinates",
		}}).
		Pipeline
}
