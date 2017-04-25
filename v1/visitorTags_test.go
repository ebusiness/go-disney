package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ebusiness/go-disney/middleware"
	"github.com/ebusiness/go-disney/v1/models"
)

// go test v1/index*.go -v
func TestVisitorIndex(t *testing.T) {
	router := gin.New()
	// mongo
	router.Use(middleware.MongoSession)
	router.GET("/test", visitorIndex)

	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// t.Log(resp.Body.String())
	if resp.Code != 200 {
		t.Fatalf("Status Error %d", resp.Code)
	}

	models := []models.VisitorTag{}

	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		t.Fatalf("Error %v", err)
	}

}
