package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

var client *mongo.Client

func Init(c *mongo.Client) {
	database := c.Database(os.Getenv("SR_EXAMPLE_MONGO_DATABASE"))
	client = c

	exampleService(database)
}

func PingDatabase() bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}

	return true
}
