package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func init() {
	control := realtimeController{}
	utils.V1.GET("/land/realtime", control.landlist)
	utils.V1.GET("/sea/realtime", control.sealist)
	utils.V1.GET("/realtime", control.list)
	utils.V1.GET("/waittimes/:id", control.waittimes)
}

type realtimeController struct{}

func (control realtimeController) landlist(c *gin.Context) {
	park := bson.M{"$match": bson.M{"base.park_kind": "1"}}
	control.search(c, park)
}

func (control realtimeController) sealist(c *gin.Context) {
	park := bson.M{"$match": bson.M{"base.park_kind": "2"}}
	control.search(c, park)
}

func (control realtimeController) list(c *gin.Context) {
	// park := bson.M{"$match": bson.M{"park_kind": bson.M{"$in": []string{"1", "2"}}}}
	control.search(c)
}

func (control realtimeController) search(c *gin.Context, bsons ...bson.M) {
	models := []models.Realtime{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			timeNow := utils.Now()
			matchToday := bson.M{
				"$match": bson.M{
					"updateTime": bson.M{
						"$gt": time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.UTC),
					},
				},
			}
			addFields := bson.M{"$addFields": bson.M{"realtime": "$$ROOT"}}
			group := bson.M{
				"$group": bson.M{
					"_id":      "$str_id",
					"realtime": bson.M{"$last": "$realtime"},
				},
			}

			pipeline = (utils.BsonCreater{}).
				Append(bson.M{"$sort": bson.M{"_id": 1, "createTime": 1, "updateTime": 1}}).
				Append(matchToday, addFields, group).
				LookupWithUnwind("attractions", "_id", "str_id", "base").
				Append(bsons...).
				Append(bson.M{
					"$addFields": bson.M{
						"base.realtime": "$realtime",
					},
				}, bson.M{"$replaceRoot": bson.M{"newRoot": "$base"}}).
				LookupWithUnwind("areas", "area", "_id", "area").
				Pipeline
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}

func (control realtimeController) waittimes(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 1 { //|| !bson.IsObjectIdHex(id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	models := []models.Waittime{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M

	utils.SafelyExecutorForGin(c,
		func() {
			timeNow := utils.Now()
			pipeline = []bson.M{
				bson.M{
					"$match": bson.M{
						"str_id": id,
						"updateTime": bson.M{
							"$gt": time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.UTC),
						},
					},
				},
				bson.M{
					"$project": bson.M{
						// "str_id":     1,
						"waitTime":   1,
						"updateTime": 1,
						"createTime": 1,
					},
				},
				bson.M{
					"$sort": bson.M{"createTime": -1},
				},
			}
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}
