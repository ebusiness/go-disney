package v1

import (
	"github.com/ahmetb/go-linq"
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"sort"
	"time"
)

func (control planController) algonrithmsWaittime(c *gin.Context, model models.PlanTemplate, datetime time.Time) models.PlanTemplate {

	mongo := middleware.GetMongo(c)

	sortList := control.sortShow(c, model, datetime)
	// log.Println("sortRoutes", sortList)

	tempRoutes := []models.PlanRoute{}
	for _, item := range model.Route {
		sorted := false
		for _, sortItem := range sortList {
			if item.StrID == sortItem.StrID {
				sorted = true
				break
			}
		}
		if !sorted {
			tempRoutes = append(tempRoutes, item)
		}
	}
	sortRoutes := []models.PlanRoute{}
	for _, sortItem := range sortList {
		if sortItem.Schedule.StartTime.Before(datetime) {
			continue
		}
		sortRoutes = append(sortRoutes, sortItem)
	}
	resultRoute := []models.PlanRoute{}
	length := len(tempRoutes) - 1
	for routeIndex, item := range tempRoutes {
		nextID := ""
		if routeIndex < length {
			nextID = tempRoutes[routeIndex+1].StrID
		}

		insertItem := item
		for sortIndex, sortItem := range sortRoutes {
			waittime := control.getPredictionWaittime(c, item, datetime)
			tempItem := control.getItemWithsSchdule(mongo, sortItem.StrID, datetime, waittime, item)
			if sortItem.Schedule.IsConflict(tempItem.Schedule) {
				insertItem = sortItem
				resultRoute = append(resultRoute, sortItem)
				datetime = sortItem.Schedule.EndTime
				sortRoutes = append(sortRoutes[:sortIndex], sortRoutes[sortIndex+1:]...)
				continue
			}
		}
		waittime := control.getPredictionWaittime(c, insertItem, datetime)
		item = control.getItemWithsSchdule(mongo, nextID, datetime, waittime, item)
		datetime = item.Schedule.EndTime

		resultRoute = append(resultRoute, item)
	}
	for _, sortItem := range sortRoutes {
		resultRoute = append(resultRoute, sortItem)
	}
	// log.Println(datetime) //end time
	model.Route = resultRoute
	return model
}

func (control planController) getItemWithsSchdule(mongo middleware.Mongo, nextID string, datetime time.Time, waittime float64, route models.PlanRoute) models.PlanRoute {
	route.Schedule.StartTime = datetime

	route.WaitTime = waittime
	route.DistanceToNext = control.getDistanceToNext(mongo, route.StrID, nextID)
	route.WalktimeToNext = math.Ceil(route.DistanceToNext / control.speed())
	cost := route.TimeCost + route.WalktimeToNext + waittime
	route.Schedule.EndTime = datetime.Add(time.Minute * time.Duration(cost))

	return route
}

func (control planController) sortShow(c *gin.Context, model models.PlanTemplate, datetime time.Time) []models.PlanRoute {
	showIds := []string{}
	showRoutes := []models.PlanRoute{}
	linq.From(model.Route).WhereT(func(item models.PlanRoute) bool {
		return item.Attraction.Category == "show"
	}).ToSlice(&showRoutes)

	linq.From(showRoutes).SelectT(func(item models.PlanRoute) string {
		return item.StrID
	}).ToSlice(&showIds)

	mongo := middleware.GetMongo(c)
	schedules := control.getScheduleList(mongo, datetime, showIds)

	routes := []models.PlanRoute{}
	linq.From(schedules).ForEachT(func(item models.ScheduleDaily) {
		route := control.getNotConflictShow(item, showRoutes, routes)
		if route == nil {
			return
		}
		routes = append(routes, *route)
	})

	sort.Sort(sortPlanRoute(routes))
	return routes
}

func (control planController) getNotConflictShow(item models.ScheduleDaily, showRoutes []models.PlanRoute, routes []models.PlanRoute) *models.PlanRoute {
	res := linq.From(item.Schedules).WhereT(func(showTime models.Schedule) bool {

		return linq.From(routes).WhereT(func(route models.PlanRoute) bool {
			return route.Schedule.IsConflict(showTime)
		}).First() == nil

	}).First()

	if res == nil {
		return nil
	}
	log.Println(res)

	schedule := res.(models.Schedule)
	route := linq.From(showRoutes).WhereT(func(showRoute models.PlanRoute) bool {
		return showRoute.StrID == item.StrID
	}).First().(models.PlanRoute)
	route.Schedule = schedule

	return &route
}
