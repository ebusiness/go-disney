package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const value = "value"

func TestRedis_SessionGetSet(t *testing.T) {
	r := gin.Default()

	r.Use(SessionRedisStore)

	r.GET("/set", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", value)
		session.Save()
		c.String(http.StatusOK, value)
	})

	r.GET("/get", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("key") != value {
			t.Error("Session writing failed")
		}
		session.Save()
		c.String(http.StatusOK, value)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/set", nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/get", nil)
	req2.Header.Set("Cookie", res1.Header().Get("Set-Cookie"))
	r.ServeHTTP(res2, req2)
}
