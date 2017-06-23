package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/algorithms"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
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
		c.AbortWithStatus(http.StatusNotAcceptable) //406
		return
	}
	t, err := time.Parse("2006-01-02 -0700", dateString+" +0900")
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable) //406
		return
	}
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
	utils.Executor(c).Parallel(utils.ParallelCallback{
		"datetime": func() (interface{}, error) {
			return datetime, nil
		},
		"prediction": func() (interface{}, error) {
			var prediction interface{}
			if datetime.After(utils.Now().AddDate(0, 0, -1)) {
				waittime := algorithms.CalculateWaitTime(control.id, datetime)
				prediction = waittime.List(c)
			}
			return prediction, nil
		},
		"realtime": func() (interface{}, error) {
			if datetime.After(utils.Now()) {
				return nil, nil
			}
			models := []models.RealWaittime{}
			mongo := middleware.GetMongo(c)
			collection := mongo.GetCollection(models)
			err := collection.Pipe(control.getPipelineOfRealtimeWaittimes(datetime)).All(&models)
			return models, err
		},
	})
}

func (control attractionController) getPipelineOfRealtimeWaittimes(datetime time.Time) []bson.M {
	nextDay := datetime.AddDate(0, 0, 1)
	return []bson.M{
		{
			"$match": bson.M{
				"str_id": control.id,
				"createTime": bson.M{
					"$gt": utils.DatetimeOfDate(datetime),
					"$lt": utils.DatetimeOfDate(nextDay),
				},
			},
		},
		{
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
		{
			"$sort": bson.M{"createTime": 1},
		},
	}
}
