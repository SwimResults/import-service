package client

import (
	"github.com/swimresults/athlete-service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExecClient() {
	team, _ := primitive.ObjectIDFromHex("64674c637f47c7e53bd06717")

	var athlete = model.Athlete{
		Name:      "Hermine Granger",
		Firstname: "Hermine",
		Lastname:  "Granger",
		Year:      2001,
		Gender:    "FEMALE",
		DsvId:     123456,
		Team: model.Team{
			Identifier: team,
		},
	}
	addAthlete(athlete)
}
