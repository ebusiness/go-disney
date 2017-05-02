package v1

import (
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"testing"
)

// go test v1/attractions*.go -v
func TestAttractions(t *testing.T) {

	control := attractionController{}
	models := []models.Attraction{}
	utils.CreaterTestForHTTP(t, "/test", "/test", control.list, &models)

	if len(models) < 1 {
		// t.Fatalf("NoData")
		// travis-ci has no data
		return
	}
	testAttractionsDetail(t, models[0].ID.Hex())
}

func testAttractionsDetail(t *testing.T, id string) {

	control := attractionController{}
	model := models.Attraction{}
	utils.CreaterTestForHTTP(t, "/test/:id", "/test/"+id, control.detail, &model)

	// t.Log(model)
}
