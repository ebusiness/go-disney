package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

func init() {
	control := scheduleController{}
	utils.V1.GET("/schedule/:datetime", control.schedule)
}

type scheduleController struct {
	baseController
}

func (control scheduleController) schedule(c *gin.Context) {
	control.initialization(c)

	getPipeline := func(param interface{}) (interface{}, error) {
		datetime, err := time.Parse("2006-01-02 -0700", c.Param("datetime")+" +0900")
		log.Println(datetime)
		return bson.M{"park": control.park, "date": datetime}, err
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.(bson.M)
		model := bson.M{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollectionByName("schedule_calendar")
		err := collection.Find(pipeline).One(&model)
		return model, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
}
