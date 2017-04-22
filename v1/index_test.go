package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// go test v1/index*.go -v
func TestIndex(t *testing.T) {
	router := gin.New()
	router.GET("/test", index)

	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// t.Log(resp.Body.String())
	if resp.Code == 200 {
		t.Log("passed")
	} else {
		t.Fatalf("Status Error %d", resp.Code)
	}
}
