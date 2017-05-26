package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
)

// Executor - execute function safely for gin
func Executor(c *gin.Context) Async {
	return Async{c}
}

// Async -
type Async struct {
	context *gin.Context
}

// ParallelCallback -
type ParallelCallback map[string]func() (interface{}, error)

// IsAborted return context.IsAborted or response is already written
func (async Async) IsAborted() bool {
	return async.context.IsAborted() || async.context.Writer.Written()
}

// defer
func (async Async) abort() {
	err := recover()
	if err == nil {
		return
	}
	log.Println("Error ", err) // stream closed

	if async.IsAborted() {
		return
	}
	async.context.AbortWithStatus(http.StatusNotAcceptable)
}

func (async Async) output(res interface{}, err error) {
	if async.IsAborted() {
		return
	}
	if err != nil || res == nil {
		async.context.AbortWithStatus(http.StatusNotFound)
		return
	}
	isSlice := reflect.TypeOf(res).Kind() == reflect.Slice
	if isSlice && reflect.ValueOf(res).Len() == 0 {
		async.context.AbortWithStatus(http.StatusNotFound)
		return
	}
	async.context.JSON(http.StatusOK, res)
}

// Waterfall - execute function safely for gin
func (async Async) Waterfall(tasks ...func(param interface{}) (interface{}, error)) {
	defer async.abort()

	var res interface{}
	var err error
	for _, task := range tasks {
		if async.IsAborted() {
			return
		}
		if err != nil { // last err should not to panic(err)->StatusNotAcceptable
			panic(err)
		}
		res, err = task(res)
	}
	async.output(res, err)
}

// Parallel - execute function safely for gin
func (async Async) Parallel(tasks ParallelCallback) {
	defer async.abort()

	type parallelResult struct {
		name string
		res  interface{}
	}

	parallelChan := make(chan parallelResult, len(tasks))

	for name := range tasks {
		nameTemp := name
		go func() {
			res, err := tasks[nameTemp]()
			if err != nil { // last err should not to panic(err)->StatusNotAcceptable
				panic(err)
			}
			parallelChan <- parallelResult{nameTemp, res}
		}()
	}
	list := map[string]interface{}{}
	for range tasks {
		temp := <-parallelChan
		if async.IsAborted() {
			return
		}
		list[temp.name] = temp.res
	}
	async.output(list, nil)
}
