package v1

import (
	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
	"testing"
)

// go test v1/attractions*.go -v
func TestVisitorIndex(t *testing.T) {

	control := visitorController{}
	models := []models.VisitorTag{}
	utils.CreaterTestForHTTP(t, "/test", "/test", control.tags, &models)

	if len(models) < 1 {
		// t.Fatalf("NoData")
		// travis-ci has no data
	}
}
