package main

import (
	"context"
	"fmt"
	client2 "github.com/swimresults/import-service/client"
	"github.com/swimresults/import-service/controller"
	"github.com/swimresults/import-service/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

var client *mongo.Client

func main() {
	//ctx := connectDB()
	service.Init(client)
	controller.Run()
	client2.ExecClient()

	//if err := client.Disconnect(ctx); err != nil {
	//	panic(err)
	//}
}

func connectDB() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	var uri = "mongodb://"
	if os.Getenv("SR_IMPORT_MONGO_USERNAME") != "" {
		uri += os.Getenv("SR_IMPORT_MONGO_USERNAME") + ":" + os.Getenv("SR_IMPORT_MONGO_PASSWORD") + "@"
	}
	uri += os.Getenv("SR_IMPORT_MONGO_HOST") + ":" + os.Getenv("SR_IMPORT_MONGO_PORT")
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_IMPORT_MONGO_HOST") + ":" + os.Getenv("SR_IMPORT_MONGO_PORT") + "' as '" + os.Getenv("SR_IMPORT_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to connect to mongo database"))
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_IMPORT_MONGO_HOST") + ":" + os.Getenv("SR_IMPORT_MONGO_PORT") + "' as '" + os.Getenv("SR_IMPORT_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to reach mongo database"))
	}

	return ctx
}
