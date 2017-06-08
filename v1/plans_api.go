package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math"
	"net/http"
)

func init() {
	control := planController{}
	utils.V1.GET("/plans", control.list)
	utils.V1.GET("/plan/random", control.random)
	utils.V1.GET("/plans/:id/:datetime", control.detail)
	utils.V1.POST("/plans/customize", control.customize)
}

type planController struct {
	baseController
}

func (control planController) list(c *gin.Context) {
	control.initialization(c)

	getPipeline := func(param interface{}) (interface{}, error) {
		return control.getConditionsForPlanList(), nil
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.([]bson.M)
		models := []models.PlanTemplate{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollection(models)
		err := collection.Pipe(pipeline).All(&models)
		return models, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
}

func (control planController) random(c *gin.Context) {
	control.initialization(c)

	getPipeline := func(param interface{}) (interface{}, error) {
		return append(control.getConditionsForPlanList(),
			bson.M{"$sample": bson.M{"size": 1}}), nil
	}
	exec := func(param interface{}) (interface{}, error) {
		pipeline := param.([]bson.M)
		model := models.PlanTemplate{}
		mongo := middleware.GetMongo(c)
		collection := mongo.GetCollection(model)
		err := collection.Pipe(pipeline).One(&model)
		return model, err
	}
	utils.Executor(c).Waterfall(getPipeline, exec)
}
func (control planController) detail(c *gin.Context) {
	control.initialization(c)
	if len(control.id) < 1 || !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	datetime := control.getDatetime(c)
	mongo := middleware.GetMongo(c)

	getPlanFromCache := func(param interface{}) (interface{}, error) {
		model, err := control.getPlan(mongo, datetime)
		if err == nil {
			return model, nil
		}
		return nil, nil
	}
	makePlan := func(param interface{}) (interface{}, error) {
		if param != nil {
			return param, nil
		}
		model := models.PlanTemplate{}
		collection := mongo.GetCollection(model)
		pipeline := control.getConditionsForPlanList(bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(control.id)}})
		err := collection.Pipe(pipeline).One(&model)
		if err != nil {
			c.JSON(http.StatusNotFound, nil)
		}
		return model, err
	}
	exec := func(param interface{}) (interface{}, error) {
		model := param.(models.PlanTemplate)
		if model.Start == nil {
			model.Start = &datetime
			for routeIndex := range model.Route {
				model.Route[routeIndex].WalktimeToNext = math.Ceil(model.Route[routeIndex].DistanceToNext / control.speed())
			}
		}
		model = control.algonrithmsWaittime(c, model, datetime)
		control.cachePlan(mongo, model)
		return model, nil
	}

	utils.Executor(c).Waterfall(getPlanFromCache, makePlan, exec)
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
		model.Route[routeIndex].WalktimeToNext = math.Ceil(distanceToNext / control.speed())
	}
	model = control.algonrithmsWaittime(c, model, *model.Start)

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
