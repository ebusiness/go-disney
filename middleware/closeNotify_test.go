package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func Test_CloseNotifyForAbort(t *testing.T) {
	r := gin.Default()

	r.Use(CloseNotify)

	r.GET("/notfound", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	rec := newCloseNotifyingRecorder()
	resp := gin.ResponseWriter(rec)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	r.ServeHTTP(resp, req)

	if resp.Status() != http.StatusNotFound {
		t.Errorf("Expected %v - Got %v", http.StatusNotFound, resp.Status())
	}
}

func Test_CloseNotifyForCancel(t *testing.T) {
	r := gin.Default()

	r.Use(CloseNotify)

	r.GET("/cancel", func(c *gin.Context) {
		time.Sleep(1000 * time.Millisecond)

		if c.IsAborted() {
			return
		}
		c.String(http.StatusOK, "hello world!")
	})

	rec := newCloseNotifyingRecorder()
	resp := gin.ResponseWriter(rec)
	req, _ := http.NewRequest("GET", "/cancel", nil)

	go func() {
		// waitting for request
		time.Sleep(100 * time.Millisecond)
		rec.close()
	}()
	r.ServeHTTP(resp, req)

	if resp.Status() != http.StatusNotAcceptable {
		t.Errorf("Expected %v - Got %v", http.StatusNotAcceptable, resp.Status())
	}
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		nil,
		make(chan bool, 1),
	}
}

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	http.Hijacker
	closed chan bool
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

func (c *closeNotifyingRecorder) Size() int {
	return 0
}

func (c *closeNotifyingRecorder) Status() int {
	return c.Code
	// return c.Result().StatusCode
}

func (c *closeNotifyingRecorder) WriteHeaderNow() {
}

func (c *closeNotifyingRecorder) Written() bool {
	return c.Body != nil
}

func (c *closeNotifyingRecorder) WriteString(s string) (n int, err error) {
	n = http.StatusOK
	err = nil
	return
}
