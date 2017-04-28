package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// test for mongo middleware
// go test middleware/mongo*.go -v
func TestMongo_ReadWrite(t *testing.T) {
	r := gin.Default()

	r.Use(MongoSession)

	r.GET("/test", func(c *gin.Context) {
		mongo := GetMongo(c)
		p := Person{}
		collection := mongo.GetCollection(p)

		testInsert(t, collection)
		testFind(t, collection)
		testDrop(t, collection)

		c.String(http.StatusOK, "hello world!")
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(res1, req1)
}

func testInsert(t *testing.T, collection *mgo.Collection) {
	// t.Log("create test_people, and insert Ale && Cla")

	err := collection.Insert(
		&Person{Name: "Ale", Phone: "+55 53 8116 9639"},
		&Person{Name: "Cla", Phone: "+55 53 8402 8510"},
	)

	if err != nil {
		t.Error(err)
	} else {
		t.Log("[testInsert] passed")
	}
}

func testFind(t *testing.T, collection *mgo.Collection) {
	// t.Log("find Ale's Phone No.")

	result := Person{}
	err := collection.Find(bson.M{"name": "Ale"}).One(&result)

	if err != nil {
		t.Error(err)
		return
	}
	if result.Phone == "+55 53 8116 9639" {
		t.Log("[testFind] passed")
	} else {
		t.Error(errors.New("the phone No. is not matched"))
	}
}

func testDrop(t *testing.T, collection *mgo.Collection) {

	err := collection.DropCollection()
	if err != nil {
		t.Error(err)
	} else {
		// t.Log("collection droped: test_people")
		t.Log("[testDrop] passed")
	}
}

type Person struct {
	collectionName string `collectionName:"test_people1"`
	Name           string
	Phone          string
}
