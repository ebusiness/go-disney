package v1

import (
	"github.com/ebusiness/go-disney/utils"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

// go test v1/schedule*.go -v
func TestScheduleIndex(t *testing.T) {

	control := scheduleController{}
	model := bson.M{}
	utils.CreaterTestForHTTP(t, "/test:datetime", "/test/2017-06-30", control.schedule, &model)
}
