package config

const (
	// A default server only if os.Getenv("MONGOSERVER") is not set
	MongoDefaultServer = "localhost"
	// A default port only if os.Getenv("MONGOPORT") is not set
	MongoDefaultPort = "27017"
	// A default database name only if os.Getenv("MONGODATABASE") is not set
	MongoDefaultDatabase = "disney"
)
