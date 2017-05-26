package algorithms

import (
	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

// CalculateWaitTime - it's will be most important in the future, now it just check if the day is holiday, weekend or specialday
func CalculateWaitTime(strID string, datetime time.Time) WaitTime {
	// log.Println("------------------------>>>>>>>>>>")
	point := 0 // 0 ~ 10000
	// if sunshine
	point += 3000

	point += calculateWaitTimeForWeekend(datetime)

	if utils.IsHoliday(datetime) {
		point += calculateWaitTimeForHoliday(datetime)
	}

	point += calculateWaitTimeForBeforeHoliday(datetime)

	point += calculateWaitTimeForSpecialDay(datetime)

	result := Normal
	if point > 6000 {
		result = Superbusy
	} else if point > 3000 {
		result = Busy
	}
	// debugPrintForCalculateWaitTime(result, point, datetime)
	// log.Println("<<<<<<<<<------------------------")
	return WaitTime{result, strID}
}

func debugPrintForCalculateWaitTime(result waitTimeRank, point int, datetime time.Time) {
	log.Println(point)
	log.Println(datetime.Date())
	log.Println(result)
}

func calculateWaitTimeForWeekend(datetime time.Time) int {
	if datetime.Weekday() == time.Sunday {
		// log.Println("Sunday")
		return 1000
	}

	if datetime.Weekday() == time.Saturday {
		// log.Println("Saturday")
		return 2000
	}
	return 0
}

func calculateWaitTimeForHoliday(datetime time.Time) int {
	if !utils.IsHoliday(datetime) {
		return 0
	}

	if datetime.Weekday() == time.Friday {
		// log.Println("(Friday) 3 days")
		return 3000
	}
	if datetime.Weekday() == time.Monday {
		// log.Println("(Monday) 3 days")
		return 1000
	}
	dayAfterTomorrow := datetime.AddDate(0, 0, 2)
	if utils.IsHoliday(dayAfterTomorrow) {
		// log.Println("(Holiday) 3 days")
		return 3000
	}

	tomorrow := datetime.AddDate(0, 0, 1)
	if utils.IsHoliday(tomorrow) {
		// log.Println("(Holiday past one day) 3 days")
		return 3000
	}
	if tomorrow.Weekday() != time.Sunday {
		// log.Println("this holiday is not Saturday - the day will eat by evil if Saturday")
		return 1000
	}

	return 0
}

func calculateWaitTimeForBeforeHoliday(datetime time.Time) int {
	tomorrow := datetime.AddDate(0, 0, 1)
	if utils.IsHoliday(tomorrow) {
		// log.Println("tomorrow is holiday too")
		return 1000
	}
	return 0
}

func calculateWaitTimeForSpecialDay(datetime time.Time) int {
	point := 0
	if datetime.Month() == time.March {
		point += 1000
		point += calculateWaitTimeForWeekend(datetime) // once more
		if datetime.Day() > 10 {
			point += 1000
		}
		if datetime.Day() > 25 {
			point += 1000
		}
	}
	if datetime.Month() == time.April {
		if datetime.Day() < 7 {
			point += 1000
		}
	}
	if datetime.Month() == time.August {
		point += 500
	}
	if datetime.Month() == time.November {
		point += calculateWaitTimeForWeekend(datetime) / 2 // once more
	}

	if datetime.Month() == time.December {
		point += calculateWaitTimeForWeekend(datetime) // once more
		if datetime.Day() == 23 || datetime.Day() == 24 || datetime.Day() == 25 {
			point += 1000
		}
		if datetime.Day() > 27 {
			point += 1500
		}

	}
	if datetime.Month() == time.January {
		if datetime.Day() < 7 {
			point += 1000
		}
	}
	return point
}

type waitTimeRank int

// rank -
const (
	Normal waitTimeRank = iota
	Busy
	Superbusy
)

var waitTimeRanks = [...]string{
	"Normal",
	"Busy",
	"Superbusy",
}

// String for print out
func (wt waitTimeRank) String() string { return waitTimeRanks[wt] }

// WaitTime -
type WaitTime struct {
	waitTimeRank
	strID string
}

// // After waitTimes after `datetime` of this day (every 15m)
// func (wt WaitTime) After(datetime time.Time, c *gin.Context) []models.PredictionWaittime {
//   log.Println(utils.Now().Weekday() == time.Thursday)
//   return wt.List(c, []bson.M{})
// }

// // Before waitTimes before `datetime` of this day (every 15m)
// func (wt WaitTime) Before(datetime time.Time) {
//   log.Println(utils.Now().Weekday() == time.Thursday)
//
// }

// List waitTime list(every 15m)
func (wt WaitTime) List(c *gin.Context) []models.PredictionWaittime {
	models := []models.PredictionWaittime{}
	mongo := middleware.GetMongo(c)
	collection := mongo.GetCollection(models)
	pipeline := wt.getConditions()

	// result := []bson.M{}
	// collection.Pipe(pipeline).All(&result)
	// log.Println(result)

	collection.Pipe(pipeline).All(&models)
	return models
}

func (wt WaitTime) getConditions() []bson.M {
	rank := "rank1"
	if wt.waitTimeRank == Superbusy {
		rank = "rank3"
	} else if wt.waitTimeRank == Busy {
		rank = "rank2"
	}

	return []bson.M{
		{"$match": bson.M{"str_id": wt.strID}},
		{
			"$project": bson.M{
				"str_id": 1,
				"rank":   "$" + rank,
			},
		},
		{"$unwind": "$rank"},
		{"$replaceRoot": bson.M{"newRoot": "$rank"}},
		{
			"$sort": bson.M{"createTime": 1},
		},
	}
}
