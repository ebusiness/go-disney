package v1

import (
	"testing"

	"github.com/ebusiness/go-disney/utils"
	"github.com/ebusiness/go-disney/v1/models"
)

// go test v1/attractions*.go -v
func TestVisitorIndex(t *testing.T) {

	models := []models.VisitorTag{}
	utils.CreaterTestForHTTP(t, "/test", "/test", visitorIndex, &models)

	if len(models) < 1 {
		// t.Fatalf("NoData")
		// travis-ci has no data
		return
	}
}
