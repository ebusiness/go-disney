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
	"math"
	"net/http"
	"time"
)

func init() {
	control := planController{}
	utils.V1.GET("/plans", control.list)
	utils.V1.GET("/plans/:id/:datetime", control.detail)
	utils.V1.POST("/plans/customize", control.customize)
}

type planController struct {
	baseController
}

func (control planController) speed(c *gin.Context) float64 {
	return 0.88 * 60 // pre minutes // old man was 0.88 pre second, young man was 1.27
}

func (control planController) list(c *gin.Context) {
	control.initialization(c)
	models := []models.PlanTemplate{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	var pipeline []bson.M
	// result := []bson.M{}
	utils.SafelyExecutorForGin(c,
		func() {
			pipeline = control.getConditionsForPlanList()
		},
		func() {
			collection.Pipe(pipeline).All(&models)
		},
		func() {
			c.JSON(http.StatusOK, models)
		})
}

func (control planController) getDatetime(c *gin.Context) time.Time {
	datetime, err := time.Parse("2006-01-02T15:04:05-0700", c.Param("datetime"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
	return datetime
}
func (control planController) detail(c *gin.Context) {
	control.initialization(c)
	if len(control.id) < 1 || !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	datetime := control.getDatetime(c)
	mongo := middleware.GetMongo(c)
	model := models.PlanTemplate{}
	var err error

	utils.SafelyExecutorForGin(c,
		func() {
			model, err = control.getPlan(mongo, datetime)
			if err == nil {
				c.JSON(http.StatusOK, control.algonrithmsWaittime(c, model, datetime))
			}
		},
		func() {
			collection := mongo.GetCollection(model)
			pipeline := control.getConditionsForPlanList(bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(control.id)}})
			err = collection.Pipe(pipeline).One(&model)
			if err != nil {
				c.JSON(http.StatusNotFound, nil)
			}
		},
		func() {
			model.Start = datetime
			for routeIndex := range model.Route {
				model.Route[routeIndex].WalktimeToNext = math.Ceil(model.Route[routeIndex].DistanceToNext / control.speed(c))
			}
			control.cachePlan(mongo, model)
			model = control.algonrithmsWaittime(c, model, datetime)
			c.JSON(http.StatusOK, model)
		})
}

func (control planController) algonrithmsWaittime(c *gin.Context, model models.PlanTemplate, datetime time.Time) models.PlanTemplate {
	for index, item := range model.Route {

		waittime := algorithms.CalculateWaitTime(item.StrID, datetime)
		predictions := waittime.List(c)

		cost := item.TimeCost + item.WalktimeToNext

		for _, prediction := range predictions {
			if datetime.Hour() > prediction.CreateTime.Hour() {
				continue
			}
			if datetime.Minute() > prediction.CreateTime.Minute() {
				continue
			}
			if prediction.WaitTime != nil {
				cost += *prediction.WaitTime
				model.Route[index].WaitTime = *prediction.WaitTime
			}
			// log.Println(item.StrID, "have to wait", datetime, prediction.CreateTime, prediction.WaitTime, item.TimeCost, item.WalktimeToNext, "=", cost)
			break
		}
		// if model.Route[index].WaitTime < 1 {
		// 	log.Println(item.StrID, "neednot to wait", datetime, item.TimeCost, item.WalktimeToNext, model.Route[index].WaitTime, "=", cost)
		// }
		datetime = datetime.Add(time.Minute * time.Duration(cost))
	}
	// log.Println(datetime) //end time
	return model
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

func (control planController) customize(c *gin.Context) {
	control.initialization(c)
	model := models.PlanTemplate{}
	err := c.BindJSON(&model)
	if err != nil {
		log.Println("customize", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	mongo := middleware.GetMongo(c)

	length := len(model.Route) - 1
	ids := []string{}
	for routeIndex := range model.Route {
		id := model.Route[routeIndex].StrID
		ids = append(ids, id)
		nextID := ""
		if routeIndex < length {
			nextID = model.Route[routeIndex+1].StrID
		}
		model.Route[routeIndex].TimeCost = control.getTimeCost(mongo, id)
		distanceToNext := control.getDistanceToNext(mongo, id, nextID)
		model.Route[routeIndex].DistanceToNext = distanceToNext
		model.Route[routeIndex].WalktimeToNext = math.Ceil(distanceToNext / control.speed(c))
	}
	model = control.algonrithmsWaittime(c, model, model.Start)

	attractions := []models.Attraction{}
	pipeline := (utils.BsonCreator{}).
		Append(bson.M{"$match": bson.M{"str_id": bson.M{"$in": ids}}}).
		Append(bson.M{"$project": bson.M{
			"str_id":           1,
			"name":             "$name." + control.lang,
			"main_visual_urls": 1,
			"category":         1,
			"is_available":     1,
		}}).
		Pipeline
	mongo.GetCollection(attractions).Pipe(pipeline).All(&attractions)
	for routeIndex := range model.Route {
		id := model.Route[routeIndex].StrID
		for _, item := range attractions {
			if id == item.StrID {
				model.Route[routeIndex].Attraction = item
				break
			}
		}
	}

	control.saveCustomizePlan(mongo, model)
	c.JSON(http.StatusOK, model)
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

func (control planController) saveCustomizePlan(mongo middleware.Mongo, template models.PlanTemplate) {
	model := models.PlanCustomize{PlanTemplate: template, Lang: control.lang}
	model.PlanTemplate.ID = bson.NewObjectId()

	collection := mongo.GetCollection(model)
	err := collection.Insert(model)
	log.Println("saveCustomizePlan", err)
}
