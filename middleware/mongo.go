package middleware

import (
	"log"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"

	"github.com/ebusiness/go-disney/config"
)

const mongoKey = "mongo"

var databaseName = func() string {
	name := os.Getenv("MONGODATABASE")
	if len(name) < 1 {
		return config.MongoDefaultDatabase
	}
	return name

}()

var session = func() *mgo.Session {
	mongoSession, err := mgo.Dial(getConnectionString())

	if err != nil {
		log.Fatalln("cannot connect to the mongo server", err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)
	return mongoSession
}()

func getConnectionString() string {
	server := os.Getenv("MONGOSERVER")
	if len(server) < 1 {
		server = config.MongoDefaultServer
	}
	port := os.Getenv("MONGOPORT")
	if len(port) < 1 {
		port = config.MongoDefaultPort
	}

	return "mongodb://" + server + ":" + port + "/" + databaseName
}

//MongoSession - MongoDB Session Storage for Connect Middleware
func MongoSession(c *gin.Context) {
	reqSession := session.Copy()
	defer reqSession.Close()
	c.Set(mongoKey, reqSession.DB(databaseName))
	c.Next()
}

//GetMongo - Get MongoDB Session store which saved by middleware
func GetMongo(c *gin.Context) Mongo {
	return Mongo{c.MustGet(mongoKey).(*mgo.Database)}
}

//Mongo - MongoDB Session store
type Mongo struct {
	database *mgo.Database
}

//GetCollection - get the Collection of model(s)
func (m *Mongo) GetCollection(i interface{}) *mgo.Collection {
	v := reflect.ValueOf(i)

	var t reflect.Type

	if v.Type().Kind() == reflect.Slice {
		t = v.Type().Elem()
	} else {
		t = v.Type()
	}

	field, found := t.FieldByName("collectionName")

	if !found {
		log.Fatalln("collection name is not set yet")
	}

	name := field.Tag.Get("collectionName")
	return m.database.C(name)
}
