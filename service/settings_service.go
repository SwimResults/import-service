package service

import (
	"context"
	"errors"
	"github.com/swimresults/import-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var settingsCollection *mongo.Collection
var settingsForMeetingNotFoundError = "no setting with given meeting found"

func settingsService(database *mongo.Database) {
	settingsCollection = database.Collection("settings")
}

func getImportSettingsByBsonDocument(d primitive.D) ([]model.ImportSetting, error) {
	var settings []model.ImportSetting

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryOptions := options.FindOptions{}

	cursor, err := settingsCollection.Find(ctx, d, &queryOptions)
	if err != nil {
		return []model.ImportSetting{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var setting model.ImportSetting
		cursor.Decode(&setting)

		settings = append(settings, setting)
	}

	if err := cursor.Err(); err != nil {
		return []model.ImportSetting{}, err
	}

	return settings, nil
}

func GetImportSettings() ([]model.ImportSetting, error) {
	return getImportSettingsByBsonDocument(bson.D{})
}

func GetImportSettingById(id primitive.ObjectID) (model.ImportSetting, error) {
	setting, err := getImportSettingsByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.ImportSetting{}, err
	}

	if len(setting) <= 0 {
		return model.ImportSetting{}, errors.New("no setting with given _id found")
	}

	return setting[0], nil
}

func GetImportSettingByMeeting(meeting string) (model.ImportSetting, error) {
	settings, err := getImportSettingsByBsonDocument(bson.D{{"meeting", meeting}})
	if err != nil {
		return model.ImportSetting{}, err
	}

	if len(settings) > 0 {
		return settings[0], nil
	}

	return model.ImportSetting{}, errors.New(settingsForMeetingNotFoundError)
}

func RemoveImportSettingById(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := settingsCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddImportSetting(setting model.ImportSetting) (model.ImportSetting, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stg, err1 := GetImportSettingByMeeting(setting.Meeting)
	if err1 != nil {
		if err1.Error() == settingsForMeetingNotFoundError {
			r, err := settingsCollection.InsertOne(ctx, setting)
			if err != nil {
				return model.ImportSetting{}, err
			}
			return GetImportSettingById(r.InsertedID.(primitive.ObjectID))
		}
		return model.ImportSetting{}, err1
	}

	_, err := settingsCollection.ReplaceOne(ctx, bson.D{{"_id", stg.Identifier}}, setting)
	if err != nil {
		return model.ImportSetting{}, err
	}

	return GetImportSettingById(stg.Identifier)

}
