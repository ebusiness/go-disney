package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/algorithms"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

func init() {
	control := attractionController{}
	utils.V1.GET("/attractions", control.list)
	utils.V1.GET("/attractions/:id", control.detail)
	utils.V1.GET("/attractions/:id/waittimes", control.waittimes)
	utils.V1.GET("/attractions/:id/waittimes/:date", control.waittimesOfDate)
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
	control.initialization(c)
	conditions := append([]bson.M{},
		bson.M{"$match": bson.M{"park_kind": control.park, "name." + control.lang: bson.M{"$ne": ""}}},
		control.commonProject(c, nil))
	control.search(c, conditions...)
}

func (control attractionController) search(c *gin.Context, conditions ...bson.M) {
	models := []models.Attraction{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			pipeline = (utils.BsonCreator{}).
				Append(conditions...).
				LookupWithUnwind("areas", "area", "_id", "area", control.lang).
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
	if len(control.id) < 1 { //|| !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	model := models.Attraction{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(model)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			basonMatchID := bson.M{"$match": bson.M{"str_id": control.id}} //bson.ObjectIdHex(control.id)}}
			project := control.commonProject(c, bson.M{"tag_ids": 1, "youtube_url": 1, "summary_tag_ids": 1, "summaries": 1})

			pipeline = (utils.BsonCreator{}).
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

func (control attractionController) waittimesOfDate(c *gin.Context) {
	dateString := c.Param("date")
	if len(dateString) < 0 {
		c.AbortWithStatus(http.StatusUnsupportedMediaType) //415
		return
	}
	t, _ := time.Parse("2006-01-02 -0700", dateString+" +0900")
	control.calculateWaittimes(c, t)
}

func (control attractionController) waittimes(c *gin.Context) {
	control.calculateWaittimes(c, utils.Now())
}

func (control attractionController) calculateWaittimes(c *gin.Context, datetime time.Time) {
	control.initialization(c)
	if len(control.id) < 1 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var prediction interface{}
	if datetime.After(utils.Now().AddDate(0, 0, -1)) {
		waittime := algorithms.CalculateWaitTime(control.id, datetime)
		prediction = waittime.List(c)
	}

	result := struct {
		Realtime   interface{} `json:"realtime,omitempty"`
		Prediction interface{} `json:"prediction,omitempty"`
	}{
		control.getRealtimeWaittimes(c, datetime),
		prediction,
	}
	c.JSON(http.StatusOK, result)
}

func (control attractionController) getRealtimeWaittimes(c *gin.Context, datetime time.Time) []models.RealWaittime {
	if datetime.After(utils.Now()) {
		log.Println("after")
		return nil
	}
	models := []models.RealWaittime{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			nextDay := datetime.AddDate(0, 0, 1)
			pipeline = []bson.M{
				bson.M{
					"$match": bson.M{
						"str_id": control.id,
						"createTime": bson.M{
							"$gt": time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, 0, time.UTC),
							"$lt": time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, time.UTC),
						},
					},
				},
				bson.M{
					"$project": bson.M{
						// "str_id":     1,
						"waitTime": 1,
						// "updateTime": 1,
						"createTime": 1,
					},
				},
				bson.M{
					"$sort": bson.M{"createTime": 1},
				},
			}
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		})
	return models
}
