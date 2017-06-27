package v1

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	// "github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os"
	"time"
)

func init() {
	control := &fileController{}
	utils.V1.GET("/files/plans/:id", control.plan)
}

type fileController struct {
	baseController
}

func (control *fileController) plan(c *gin.Context) {
	control.initialization(c)

	if len(control.id) < 1 || !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	filepath := "/asset/plans/" + control.park + control.id
	// if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
	// 	// path/to/whatever does not exist
	// }
	if _, err := os.Stat(filepath); err != nil {

		mongo := middleware.GetMongo(c)
		planData := control.getPlan(mongo)
		if len(planData) == 0 {
			planData = control.getPlanTemplate(mongo)
		}
		control.drawImage(filepath, planData)
	}
	c.File(filepath)
}

func (control fileController) drawImage(filepath string, planData []bson.M) {
	if len(planData) == 0 {
		return
	}
	resID := "land"
	if "2" == control.park {
		resID = "sea"
	}
	draw := utils.NewPlanDraw(resID)
	points := []utils.DrawPoint{}
	for _, item := range planData {
		coord := item["coord"].(bson.M)
		point := utils.DrawPoint{coord["x"].(float64), coord["y"].(float64)}
		points = append(points, point)
	}
	// draw.DrawMark("showwait20", utils.DrawPoint{200, 200})
	draw.DrawLines(points)

	for _, item := range planData {
		category := item["category"]
		coord := item["coord"].(bson.M)
		point := utils.DrawPoint{coord["x"].(float64), coord["y"].(float64)}

		rank := control.getRankString(item["waitTime"])
		markName := "wait" + rank
		if "show" == category {
			markName = "show" + markName
		}

		draw.DrawMark(markName, point)
	}
	draw.SaveImage(filepath)
	log.Println(filepath)
	// log.Println(planData)
}

func (control fileController) getRankString(elem interface{}) (rank string) {
	rank = "20"
	if elem == nil {
		return
	}
	waitTime := elem.(float64)
	if waitTime > 70 {
		rank = "70"
	} else if waitTime > 60 {
		rank = "60"
	} else if waitTime > 50 {
		rank = "50"
	} else if waitTime > 40 {
		rank = "40"
	} else if waitTime > 30 {
		rank = "30"
	}
	log.Println(waitTime, rank)
	return
}

// copy from planController
func (control fileController) getDatetime(c *gin.Context) time.Time {
	datetime, err := time.Parse("2006-01-02T15:04:05-0700", c.Param("datetime"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
	return datetime
}

func (control fileController) makePlanFiles() {

}

func (control fileController) getPlan(mongo middleware.Mongo) []bson.M {
	match := bson.M{
		"_id": bson.ObjectIdHex(control.id),
	}
	pipeline := []bson.M{
		{"$match": match},
		{"$unwind": "$route"},
		{
			"$lookup": bson.M{
				"from":         "attractions_locations",
				"localField":   "route.str_id",
				"foreignField": "str_id",
				"as":           "location",
			},
		},
		{
			"$project": bson.M{
				"str_id":   "$route.str_id",
				"category": "$route.attraction.category",
				"waitTime": "$route.waitTime",
				"coord":    "$location.point",
			},
		},
		{"$unwind": "$coord"},
	}
	log.Println(pipeline)

	res := []bson.M{}
	mongo.GetCollectionByName("cache_plans").Pipe(pipeline).All(&res)
	return res
}

func (control fileController) getPlanTemplate(mongo middleware.Mongo) []bson.M {
	match := bson.M{
		"_id": bson.ObjectIdHex(control.id),
	}
	pipeline := []bson.M{
		{"$match": match},
		{"$unwind": "$route"},
		{
			"$lookup": bson.M{
				"from":         "attractions_locations",
				"localField":   "route.str_id",
				"foreignField": "str_id",
				"as":           "location",
			},
		},
		{
			"$project": bson.M{
				"str_id": "$route.str_id",
				"coord":  "$location.point",
			},
		},
		{"$unwind": "$coord"},
	}

	res := []bson.M{}
	collection := mongo.GetCollectionByName("plan_templates")
	collection.Pipe(pipeline).
		All(&res)
	return res
}
