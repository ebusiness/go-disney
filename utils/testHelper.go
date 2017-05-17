package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ebusiness/go-disney/middleware"
)

// CreaterTestForHTTP for httptest, return statusCode
func CreaterTestForHTTP(t *testing.T, routeURL, requestURL string, handler gin.HandlerFunc, models interface{}) int {
	router := gin.New()
	// mongo
	router.Use(middleware.MongoSession)

	router.GET(routeURL, handler)

	req, _ := http.NewRequest("GET", requestURL, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK && resp.Code != http.StatusNotFound {
		t.Fatalf("Error: url (%v) status (%d)", requestURL, resp.Code)
	}

	if models == nil {
		return resp.Code
	}

	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		t.Fatalf("Error: %v", err)
	}
	return resp.Code
	// t.Log(resp.Body.String())
	// t.Log(models)
}
