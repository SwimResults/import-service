package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var client *mongo.Client

func Init(c *mongo.Client) {
	//database := c.Database(os.Getenv("SR_IMPORT_MONGO_DATABASE"))
	//client = c

	//exampleService(database)
}

func PingDatabase() bool {

	if client == nil {
		return true
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}

	return true
}
