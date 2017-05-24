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
	utils.V1.GET("/attractions/:id/waittimes", control.waittimes)
	utils.V1.GET("/attractions/:id/waittimes/:date", control.waittimesOfDate)
}

func (control attractionController) waittimesOfDate(c *gin.Context) {
	dateString := c.Param("date")
	if len(dateString) < 0 {
		c.AbortWithStatus(http.StatusNotAcceptable) //415
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
		Datetime   time.Time   `json:"datetime"`
		Realtime   interface{} `json:"realtime,omitempty"`
		Prediction interface{} `json:"prediction,omitempty"`
	}{
		datetime,
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
							"$gt": utils.DatetimeOfDate(datetime),
							"$lt": utils.DatetimeOfDate(nextDay),
						},
					},
				},
				bson.M{
					"$project": bson.M{
						// "str_id":     1,
						"waitTime": 1,
						// "updateTime": 1,
						"createTime": 1,
						"available":  1,
						// "operation_start": 1,
						"operation_end": 1,
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
