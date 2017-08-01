package v1

import (
	"github.com/ahmetb/go-linq"
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	// "log"
	"math"
	"sort"
	"time"
)

func (control planController) algonrithmsWaittime(c *gin.Context, model models.PlanTemplate, datetime time.Time) models.PlanTemplate {

	mongo := middleware.GetMongo(c)
	fpList := control.fpList(c, model, datetime)
	sortList := control.sortShow(c, model, datetime, fpList)

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
			item.FastPass = nil //invalid fp
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

		for sortIndex, sortItem := range sortRoutes {
			waittime := control.getPredictionWaittime(c, item, datetime)
			tempItem := control.getItemWithsSchdule(mongo, sortItem.StrID, datetime, waittime, item)
			if sortItem.Schedule.IsConflict(tempItem.Schedule) {

				waittime := control.getPredictionWaittime(c, sortItem, datetime)

				//　TODO　it's should be next show's StrID or `item.StrID`
				// set it to `item.StrID` temporary
				insertItem := control.getItemWithsSchdule(mongo, item.StrID, datetime, waittime, sortItem)
				if insertItem != nil {
					resultRoute = append(resultRoute, *insertItem)
					datetime = insertItem.Schedule.EndTime
				}
				if len(sortRoutes) == 1 {
					sortRoutes = []models.PlanRoute{}
				} else {
					sortRoutes = append(sortRoutes[:sortIndex], sortRoutes[sortIndex+1:]...)
				}
				continue
			}
		}
		if len(nextID) == 0 && len(sortRoutes) > 0 {
			nextID = sortRoutes[0].StrID
		}
		waittime := control.getPredictionWaittime(c, item, datetime)
		item = *control.getItemWithsSchdule(mongo, nextID, datetime, waittime, item)
		datetime = item.Schedule.EndTime

		resultRoute = append(resultRoute, item)
	}

	length = len(sortRoutes) - 1
	for index, sortItem := range sortRoutes {

		nextID := ""
		if index < length {
			nextID = sortRoutes[index+1].StrID
		}
		waittime := control.getPredictionWaittime(c, sortItem, datetime)
		insertItem := control.getItemWithsSchdule(mongo, nextID, datetime, waittime, sortItem)
		if insertItem != nil {
			resultRoute = append(resultRoute, *insertItem)
			datetime = insertItem.Schedule.EndTime
		}
	}
	// log.Println(datetime) //end time
	model.Route = resultRoute
	return model
}

func (control planController) getItemWithsSchdule(mongo middleware.Mongo, nextID string, datetime time.Time, waittime float64, route models.PlanRoute) *models.PlanRoute {
	if route.Schedule.StartTime.Before(datetime) {
		route.Schedule.StartTime = datetime
	}

	route.WaitTime = waittime
	if route.FastPass != nil {
		datetime.Before(*route.FastPass.Begin)
		if datetime.Before(*route.FastPass.Begin) {
			route.Schedule.StartTime = *route.FastPass.Begin
		}
		route.WaitTime = waittime * control.fastRate()
	}
	route.DistanceToNext = control.getDistanceToNext(mongo, route.StrID, nextID)
	route.WalktimeToNext = math.Ceil(route.DistanceToNext / control.speed())
	cost := route.TimeCost + route.WalktimeToNext + route.WaitTime
	route.Schedule.EndTime = route.Schedule.StartTime.Add(time.Minute * time.Duration(cost))

	if nil != route.FastPass && route.Schedule.EndTime.After(*route.FastPass.End) {
		return nil
	}

	return &route
}

func (control planController) fpList(c *gin.Context, model models.PlanTemplate, datetime time.Time) []models.PlanRoute {
	routes := []models.PlanRoute{}
	linq.From(model.Route).WhereT(func(item models.PlanRoute) bool {
		return item.FastPass != nil && datetime.Before(*item.FastPass.End)
	}).SelectT(func(item models.PlanRoute) models.PlanRoute {
		item.Schedule.StartTime = *item.FastPass.Begin
		if datetime.After(item.Schedule.StartTime) {
			item.Schedule.StartTime = datetime
		}
		item.Schedule.EndTime = *item.FastPass.End
		return item
	}).ToSlice(&routes)
	return routes
}

func (control planController) sortShow(c *gin.Context, model models.PlanTemplate, datetime time.Time, routes []models.PlanRoute) []models.PlanRoute {
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
	// log.Println(schedules)

	linq.From(schedules).ForEachT(func(item models.ScheduleDaily) {
		// log.Println(item)
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
	// log.Println(res)

	schedule := res.(models.Schedule)
	route := linq.From(showRoutes).WhereT(func(showRoute models.PlanRoute) bool {
		return showRoute.StrID == item.StrID
	}).First().(models.PlanRoute)
	route.Schedule = schedule

	return &route
}
