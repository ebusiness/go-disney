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
	"strconv"
	"time"
)

func init() {
	control := &fileController{}
	utils.V1.GET("/files/plans/:id", control.planIndex)
	utils.V1.GET("/files/plans/:id/waittime", control.planWaittime)
	utils.V1.GET("/files/plans/:id/categroy", control.planCategroy)
}

type fileController struct {
	baseController
}

func (control *fileController) planIndex(c *gin.Context) {
	control.drawImage(c, "index", true, func(draw *utils.PlanDraw, index int, item bson.M, point utils.DrawPoint) {
		draw.DrawMark("mark20", point)

		draw.DrawMark("round", point.Add(0, -42))

		var x float64 = -7
		if index > 8 {
			x = -13.0
		}
		draw.DrawString(strconv.Itoa(index+1), point.Add(x, -52))
	})
}
func (control *fileController) planWaittime(c *gin.Context) {
	control.drawImage(c, "waittime", false, func(draw *utils.PlanDraw, index int, item bson.M, point utils.DrawPoint) {
		rank := control.getRankString(item["waitTime"])

		draw.DrawMark("mark"+rank, point)

		draw.DrawMark("waiticon"+rank, point.Add(0, -42))
	})
}
func (control *fileController) planCategroy(c *gin.Context) { //with FP
	control.drawImage(c, "waittime", false, func(draw *utils.PlanDraw, index int, item bson.M, point utils.DrawPoint) {
		rank := control.getRankString(item["waitTime"])

		draw.DrawMark("mark"+rank, point)

		category := item["category"]
		if category != nil {
			draw.DrawMark(category.(string), point.Add(0, -42))
		}

		fastpass := item["fastpass"]
		if fastpass != nil {
			draw.DrawMark("fp", point.Add(26, -68))
		}
	})
}

func (control *fileController) drawImage(c *gin.Context, filename string,
	isDrawLines bool,
	drawFunc func(*utils.PlanDraw, int, bson.M, utils.DrawPoint)) {

	filepath := control.beforeDraw(c, filename)

	if len(filepath) == 0 {
		return
	}

	mongo := middleware.GetMongo(c)
	planData := control.getPlan(mongo)
	if len(planData) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	resID := "land"
	if "2" == control.park {
		resID = "sea"
	}
	withbackground := c.Query("withbg")

	draw := utils.NewPlanDraw(resID, withbackground != "")

	if isDrawLines {
		control.drawLines(draw, planData)
	}

	for index, item := range planData {
		point := control.getPoint(item["coord"])
		drawFunc(draw, index, item, point)
	}
	draw.SaveImage(filepath)
	c.File(filepath)
}

func (control *fileController) beforeDraw(c *gin.Context, filename string) string {

	control.initialization(c)

	if len(control.id) < 1 || !bson.IsObjectIdHex(control.id) {
		c.AbortWithStatus(http.StatusNotFound)
		return ""
	}
	path := "/asset/plans/" + control.park + "_" + control.id
	os.MkdirAll(path, 0777)
	filepath := path + "/" + filename
	// if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
	// 	// path/to/whatever does not exist
	// }
	if _, err := os.Stat(filepath); err == nil && false {
		c.File(filepath)
		return ""
	}
	return filepath
}

func (control *fileController) drawLines(draw *utils.PlanDraw, planData []bson.M) {
	points := []utils.DrawPoint{}
	for _, item := range planData {
		points = append(points, control.getPoint(item["coord"]))
	}
	// draw.DrawMark("showwait20", utils.DrawPoint{200, 200})
	draw.DrawLines(points)
}

func (control fileController) getPoint(elem interface{}) utils.DrawPoint {
	coord := elem.(bson.M)
	y := coord["y"].(float64)
	if y < 96 {
		y = 96
	}
	x := coord["x"].(float64)
	if x > 700 {
		x -= 48
	}
	return utils.DrawPoint{X: x, Y: y}
}
func (control fileController) getRankString(elem interface{}) (rank string) {
	rank = "20"
	if elem == nil {
		return
	}
	waitTime := elem.(float64) * 4
	if waitTime > 60 {
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

func (control fileController) getPlan(mongo middleware.Mongo) []bson.M {
	match := bson.M{
		"_id": bson.ObjectIdHex(control.id),
		// "$or": []bson.M{
		// 	{"template_id": bson.ObjectIdHex(control.id)},
		// 	{"_id": bson.ObjectIdHex(control.id)},
		// },
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
				"fastpass": "$route.fastpass",
				"coord":    "$location.point",
			},
		},
		{"$unwind": "$coord"},
	}
	log.Println(pipeline)

	res := []bson.M{}
	mongo.GetCollectionByName("cache_plans").Pipe(pipeline).All(&res)
	if len(res) == 0 {
		res = control.getPlanTemplate(mongo)
	}
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
			"$lookup": bson.M{
				"from":         "attractions",
				"localField":   "route.str_id",
				"foreignField": "str_id",
				"as":           "attraction",
			},
		},
		{
			"$project": bson.M{
				"str_id":   "$route.str_id",
				"category": "$attraction.category",
				"coord":    "$location.point",
			},
		},
		{"$unwind": "$coord"},
		{"$unwind": "$category"},
	}

	res := []bson.M{}
	collection := mongo.GetCollectionByName("plan_templates")
	collection.Pipe(pipeline).
		All(&res)
	return res
}
