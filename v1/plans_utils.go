package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/algorithms"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

//sort function
type sortPlanRoute []models.PlanRoute

func (slice sortPlanRoute) Len() int {
	return len(slice)
}

func (slice sortPlanRoute) Less(i, j int) bool {
	return slice[i].Schedule.StartTime.Before(slice[j].Schedule.StartTime)
}

func (slice sortPlanRoute) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//
// planController
//

func (control planController) speed() float64 {
	return 0.88 * 60 // pre minutes // old man was 0.88 pre second, young man was 1.27
}

func (control planController) getDatetime(c *gin.Context) time.Time {
	datetime, err := time.Parse("2006-01-02T15:04:05-0700", c.Param("datetime"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
	return datetime
}

func (control planController) getTimeCost(mongo middleware.Mongo, id string) float64 {
	model := models.TimeCost{}
	collection := mongo.GetCollection(model)
	collection.Find(bson.M{"str_id": id}).One(&model)
	if model.Cost == nil {
		return 0
	}
	return *model.Cost
}

func (control planController) getDistanceToNext(mongo middleware.Mongo, id, next string) float64 {
	if len(next) == 0 {
		return 0
	}
	model := models.AttractionRelations{}
	collection := mongo.GetCollection(model)
	collection.Find(bson.M{"from": id, "to": next}).One(&model)
	if model.Distance == nil {
		return 0
	}
	return *model.Distance
}

func (control planController) getScheduleList(mongo middleware.Mongo, datetime time.Time, showIds []string) []models.ScheduleDaily {
	schedules := []models.ScheduleDaily{}
	collection := mongo.GetCollection(schedules)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"startTime": bson.M{
					"$gt": utils.DatetimeOfDate(datetime),
					"$lt": utils.DatetimeOfDate(datetime.AddDate(0, 0, 1)),
				},
				"str_id": bson.M{
					"$in": showIds,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":    "$str_id",
				"str_id": bson.M{"$first": "$str_id"},
				"schedules": bson.M{"$push": bson.M{
					"startTime": "$startTime",
					"endTime":   "$endTime",
				}},
			},
		},
	}
	collection.Pipe(pipeline).All(&schedules)
	return schedules
}

func (control planController) getConditionsForPlanList(conditions ...bson.M) []bson.M {
	return (utils.BsonCreator{}).
		Append(conditions...).
		Append(bson.M{"$addFields": bson.M{"name": "$name." + control.lang, "introduction": "$introduction." + control.lang}}).
		Append(bson.M{"$addFields": bson.M{"old": "$$ROOT"}}).
		Append(bson.M{"$unwind": "$route"}).
		LookupWithUnwind("attractions", "route.str_id", "str_id", "temp", "").
		Append(bson.M{"$match": bson.M{"temp.park_kind": control.park}}).
		// LookupWithUnwind("waittimes", "route.str_id", "str_id", "waittimes", "").
		Append(bson.M{"$unwind": "$temp"}).
		// Append(bson.M{"$unwind": "$waittimes"}).
		Append(bson.M{"$project": bson.M{
			"name":  1,
			"old":   1,
			"route": 1,
			"attraction": bson.M{
				"name":             "$temp.name." + control.lang,
				"main_visual_urls": "$temp.main_visual_urls",
				"category":         "$temp.category",
				"is_available":     "$temp.is_available",
				// "waittimes":        "$waittimes",
			},
		}}).
		Append(bson.M{"$addFields": bson.M{"route.attraction": "$attraction"}}).
		Append(bson.M{"$group": bson.M{"_id": "$name", "old": bson.M{"$first": "$old"}, "route": bson.M{"$push": "$route"}}}).
		Append(bson.M{"$addFields": bson.M{"old.route": "$route"}}).
		Append(bson.M{"$replaceRoot": bson.M{"newRoot": "$old"}}).
		Pipeline
}

func (control planController) getPredictionWaittime(c *gin.Context, item models.PlanRoute, datetime time.Time) float64 {
	waittime := algorithms.CalculateWaitTime(item.StrID, datetime)
	predictions := waittime.List(c)
	for _, prediction := range predictions {
		if datetime.Hour() > prediction.CreateTime.Hour() {
			continue
		}
		if datetime.Minute() > prediction.CreateTime.Minute() {
			continue
		}
		if prediction.WaitTime != nil {
			return *prediction.WaitTime
		}
	}
	return 0
}

func (control planController) getPlan(mongo middleware.Mongo, datetime time.Time) (models.PlanTemplate, error) {
	model := models.Plan{}
	collection := mongo.GetCollection(model)
	err := collection.Find(bson.M{"template_id": bson.ObjectIdHex(control.id), "start": datetime, "lang": control.lang}).
		One(&model)
	return model.PlanTemplate, err
}

func (control planController) cachePlan(mongo middleware.Mongo, template models.PlanTemplate) {
	model := models.Plan{PlanTemplate: template, Lang: control.lang}
	model.TemplateID = model.PlanTemplate.ID
	model.PlanTemplate.ID = bson.NewObjectId()
	collection := mongo.GetCollection(model)
	control.createIndex(collection)
	collection.Insert(model)
}

func (control planController) createIndex(collection *mgo.Collection) {
	index := mgo.Index{
		Key: []string{"template_id", "start", "lang"},
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		log.Println("EnsureIndex", err)
	}
}

func (control planController) saveCustomizePlan(mongo middleware.Mongo, template models.PlanTemplate) {
	model := models.PlanCustomize{PlanTemplate: template, Lang: control.lang}
	model.PlanTemplate.ID = bson.NewObjectId()

	collection := mongo.GetCollection(model)
	err := collection.Insert(model)
	log.Println("saveCustomizePlan", err)
}
