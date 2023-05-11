package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Example struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content    string             `json:"content,omitempty" bson:"content,omitempty"`
	Number     int                `json:"number,omitempty" bson:"number,omitempty"`
}
