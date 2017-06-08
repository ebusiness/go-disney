package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
	"time"
)

func (control planController) algonrithmsWaittime(c *gin.Context, model models.PlanTemplate, datetime time.Time) models.PlanTemplate {

	mongo := middleware.GetMongo(c)

	sortList := control.sortShow(c, model, datetime)
	// log.Println("sortRoutes", sortRoutes)

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
		item = control.getItemWithsSchdule(mongo, nextID, datetime, waittime, insertItem)
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
	for _, item := range model.Route {
		if item.Attraction.Category != "show" {
			continue
		}
		showIds = append(showIds, item.StrID)
	}

	mongo := middleware.GetMongo(c)
	schedules := control.getScheduleList(mongo, datetime, showIds)

	routes := []models.PlanRoute{}

	lists := []models.Schedule{}
	for index, route := range model.Route {
		if route.Attraction.Category != "show" {
			continue
		}

		for _, item := range schedules {
			if item.StrID != model.Route[index].StrID {
				continue
			}
			// model.Route[index].Schedule = item.Schedules[0]
			for _, showTime := range item.Schedules {
				// if route.TimeCost < showTime.EndTime.Sub(*showTime.StartTime).Minutes() {
				// 	routes = append(routes, route)
				// 	break
				// }
				conflict := false
				for _, cacheTime := range lists {
					if cacheTime.IsConflict(showTime) {
						conflict = true
						break
					}
				}
				if !conflict {
					route.Schedule = showTime
					lists = append(lists, showTime)
					routes = append(routes, route)
					break
				}
			}
		}
	}
	sort.Sort(sortPlanRoute(routes))
	return routes

}

// func (control planController) needInsertShowRightNow(schedules []models.ScheduleDaily, thisRoute, nextRoute models.PlanRoute, datetime time.Time) bool {
//
// }
