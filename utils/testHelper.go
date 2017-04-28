package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ebusiness/go-disney/middleware"
)

// CreaterTestForHTTP for httptest
func CreaterTestForHTTP(t *testing.T, routeURL, requestURL string, handler gin.HandlerFunc, models interface{}) {
	router := gin.New()
	// mongo
	router.Use(middleware.MongoSession)

	router.GET(routeURL, handler)

	req, _ := http.NewRequest("GET", requestURL, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// t.Log(resp.Body.String())
	if resp.Code != http.StatusOK {
		t.Fatalf("Error: url (%v) status (%d)", requestURL, resp.Code)
	}

	if models == nil {
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		t.Fatalf("Error: %v", err)
	}
	// t.Log(models)
}
