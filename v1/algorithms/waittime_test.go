package algorithms

import (
	"github.com/ebusiness/go-disney/utils"
	"testing"
	"time"
)

func TestWaittimes(t *testing.T) {
	start := time.Date(2011, time.November, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now() //.In(jst)

	// reslut := make([]interface{}, 0)
	for start.Before(end) {
		waittime := CalculateWaitTime("", start)
		// item := []string{start.String(), waittime.String()}
		// reslut = append(reslut, item)
		// t.Log(start.String(), waittime.String())
		if start.Weekday() == time.Sunday || start.Weekday() == time.Saturday || utils.IsHoliday(start) {
			if waittime.waitTimeRank == Normal {
				t.Fatal("the day[" + start.String() + "] should not be [Normal]")
			}
		}
		start = start.AddDate(0, 0, 1)
	}
}
