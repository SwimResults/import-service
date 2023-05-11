package service

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sr-example/example-service/model"
	"time"
)

var collection *mongo.Collection

func exampleService(database *mongo.Database) {
	collection = database.Collection("example")
}

func GetExamples() ([]model.Example, error) {
	var examples []model.Example

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return []model.Example{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var example model.Example
		cursor.Decode(&example)
		examples = append(examples, example)
	}

	if err := cursor.Err(); err != nil {
		return []model.Example{}, err
	}

	return examples, nil
}

func GetExampleById(id primitive.ObjectID) (model.Example, error) {
	var example model.Example

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{{"_id", id}})
	if err != nil {
		return model.Example{}, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		cursor.Decode(&example)
		return example, nil
	}

	return model.Example{}, errors.New("no entry with given id found")
}

func RemoveExampleById(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddExample(example model.Example) (model.Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := collection.InsertOne(ctx, example)
	if err != nil {
		return model.Example{}, err
	}

	return GetExampleById(r.InsertedID.(primitive.ObjectID))
}

func UpdateExample(example model.Example) (model.Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.D{{"_id", example.Identifier}}, example)
	if err != nil {
		return model.Example{}, err
	}

	return GetExampleById(example.Identifier)
}
