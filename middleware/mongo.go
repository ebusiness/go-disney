package middleware

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"

	"github.com/ebusiness/go-disney/config"
)

const sessionKey = "mongo"

//MongoSession - MongoDB Session Storage for Connect Middleware
func MongoSession() gin.HandlerFunc {
	connString, database := getMongoServerAddress()
	session, err := mgo.Dial(connString)
	if err != nil {
		log.Fatalln("cannot connect to the mongo server", err)
	}
	session.SetMode(mgo.Monotonic, true)

	return func(c *gin.Context) {
		reqSession := session.Copy()
		defer reqSession.Close()
		c.Set(sessionKey, reqSession.DB(database))

		c.Next()
	}
}

//GetMongoSession - Get MongoDB Session Storage which saved by middleware
func GetMongoSession(c *gin.Context) *mgo.Database {
	return c.MustGet(sessionKey).(*mgo.Database)
}

func getMongoServerAddress() (connString, database string) {
	server := os.Getenv("MONGOSERVER")
	if len(server) < 1 {
		server = config.MongoDefaultServer
	}
	port := os.Getenv("MONGOPORT")
	if len(port) < 1 {
		port = config.MongoDefaultPort
	}
	database = os.Getenv("MONGODATABASE")
	if len(database) < 1 {
		database = config.MongoDefaultDatabase
	}
	connString = "mongodb://" + server + ":" + port + "/" + database
	return
}
