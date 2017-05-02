package v1

import (
	"testing"

	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
)

// go test v1/attractions*.go -v
func TestRealtime(t *testing.T) {

	control := realtimeController{}
	models := []models.Realtime{}
	utils.CreaterTestForHTTP(t, "/test", "/test", control.list, &models)

	if len(models) < 1 {
		// t.Fatalf("NoData")
		// travis-ci has no data
		return
	}
  // t.Log(models[0].StrID)
	testWaittime(t, models[0].StrID)
}

func testWaittime(t *testing.T, id string) {

	control := realtimeController{}
	model := []models.Waittime{}
	utils.CreaterTestForHTTP(t, "/test/:id", "/test/"+id, control.waittimes, &model)

	// t.Log(model)
}
