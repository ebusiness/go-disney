package v1

import (
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"testing"
)

// go test v1/plans*.go -v
func TestPlans(t *testing.T){

	control := planController{}
	models := []models.PlanTemplate{}
	utils.CreaterTestForHTTP(t, "/test", "/test", control.list, &models)

	if len(models) < 1 {
		// t.Fatalf("NoData")
		// travis-ci has no data
		return
	}
	testPlansDetail(t, models[0].ID.Hex())
}

func testPlansDetail(t *testing.T, id string) {

	control := planController{}
	model := models.PlanTemplate{}
	utils.CreaterTestForHTTP(t, "/test/:id/:datetime", "/test/"+id+"/2017-05-03T09:00:00+0900", control.detail, &model)

	// t.Log(model)
}
